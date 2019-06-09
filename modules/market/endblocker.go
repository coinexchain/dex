package market

import (
	"crypto/sha256"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/match"
)

// Some handlers which are useful when orders are matched and traded.
type InfoForDeal struct {
	bxKeeper      ExpectedBankxKeeper
	dataHash      []byte
	changedOrders map[string]*Order
	lastPrice     sdk.Dec
	context       sdk.Context
}

// returns true when a buyer's frozen money is not enough to buy LeftStock.
func notEnoughMoney(order *Order) bool {
	return order.Side == match.BUY &&
		order.Freeze < order.Price.Mul(sdk.NewDec(order.LeftStock)).RoundInt64()
}

// Wrapper an order with InfoForDeal, which contains useful handlers
type WrappedOrder struct {
	order       *Order
	infoForDeal *InfoForDeal
}

// WrappedOrder implements OrderForTrade interface
func (wo *WrappedOrder) GetPrice() sdk.Dec {
	return wo.order.Price
}

func (wo *WrappedOrder) GetAmount() int64 {
	if notEnoughMoney(wo.order) {
		return 0
	}
	return wo.order.LeftStock
}

func (wo *WrappedOrder) GetHeight() int64 {
	return wo.order.Height
}

func (wo *WrappedOrder) GetSide() int {
	return int(wo.order.Side)
}

func (wo *WrappedOrder) GetOwner() match.Account {
	return wo.order.Sender
}

func (wo *WrappedOrder) String() string {
	return wo.order.OrderID()
}

func (wo *WrappedOrder) GetHash() []byte {
	res := sha256.Sum256(append([]byte(wo.order.OrderID()), wo.infoForDeal.dataHash...))
	return res[:]
}

func (wo *WrappedOrder) Deal(otherSide match.OrderForTrade, amount int64, price sdk.Dec) {
	other := otherSide.(*WrappedOrder)
	buyer, seller := wo.order, other.order
	if buyer.Side == match.SELL {
		buyer, seller = other.order, wo.order
	}
	stockAndMoney := strings.Split(buyer.Symbol, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	// buyer and seller will exchange stockCoins and moneyCoins
	stockCoins := sdk.Coins{sdk.NewCoin(stock, sdk.NewInt(amount))}
	moneyAmount := price.MulInt(sdk.NewInt(amount)).RoundInt64()
	moneyCoins := sdk.Coins{sdk.NewCoin(money, sdk.NewInt(moneyAmount))}
	//fmt.Printf("here price:%s stock:%d money:%d seller.Freeze:%d buyer.Freeze:%d %d %s\n",
	//price.String(), amount, moneyAmount, seller.Freeze, buyer.Freeze, buyer.Quantity, buyer.Price.String())

	buyer.LeftStock -= amount
	seller.LeftStock -= amount
	buyer.Freeze -= moneyAmount
	seller.Freeze -= amount
	buyer.DealStock += amount
	seller.DealStock += amount
	buyer.DealMoney += moneyAmount
	seller.DealMoney += moneyAmount
	ctx := wo.infoForDeal.context
	// exchange the coins
	wo.infoForDeal.bxKeeper.UnFreezeCoins(ctx, seller.Sender, stockCoins)
	wo.infoForDeal.bxKeeper.SendCoins(ctx, seller.Sender, buyer.Sender, stockCoins)
	wo.infoForDeal.bxKeeper.UnFreezeCoins(ctx, buyer.Sender, moneyCoins)
	wo.infoForDeal.bxKeeper.SendCoins(ctx, buyer.Sender, seller.Sender, moneyCoins)

	// record the changed orders for further processing
	wo.infoForDeal.changedOrders[buyer.OrderID()] = buyer
	wo.infoForDeal.changedOrders[seller.OrderID()] = seller

	// record the last executed price, which will be stored in MarketInfo
	wo.infoForDeal.lastPrice = price
}

// unfreeze an ask order's stock or a bid order's money
func unfreezeCoinsForOrder(ctx sdk.Context, bxKeeper ExpectedBankxKeeper, order *Order) {
	stockAndMoney := strings.Split(order.Symbol, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	frozenToken := stock
	if order.Side == match.BUY {
		frozenToken = money
	}
	coins := sdk.Coins([]sdk.Coin{sdk.NewCoin(frozenToken, sdk.NewInt(order.Freeze))})
	bxKeeper.UnFreezeCoins(ctx, order.Sender, coins)
	if order.FrozenFee != 0 {
		coins = sdk.Coins([]sdk.Coin{sdk.NewCoin("cet", sdk.NewInt(order.FrozenFee))})
		bxKeeper.UnFreezeCoins(ctx, order.Sender, coins)
		//actualFee:=sdk.NewDec(order.DealStock).Mul(sdk.NewDec(order.FrozenFee)).Quo(sdk.NewDec(order.Quantity))
		//TODO deduct actualFee from order.Sender
	}
}

// remove the orders whose age are older than height
func removeOrderOlderThan(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, height int64) {
	for _, order := range orderKeeper.GetOlderThan(ctx, height) {
		removeOrder(ctx, orderKeeper, bxKeeper, order)
	}
}

// unfreeze the frozen token in the order and remove it from the market
func removeOrder(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, order *Order) {
	if order.Freeze != 0 {
		unfreezeCoinsForOrder(ctx, bxKeeper, order)
	}
	orderKeeper.Remove(ctx, order)
}

// Iterate the candidate orders for matching, and remove the orders whose sender is forbidden by the money owner or the stock owner.
func filterCandidates(ctx sdk.Context, asKeeper ExpectedAssertStatusKeeper, ordersIn []*Order, stock, money string) []*Order {
	ordersOut := make([]*Order, 0, len(ordersIn))
	for _, order := range ordersIn {
		if !(asKeeper.IsForbiddenByTokenIssuer(ctx, stock, order.Sender) ||
			asKeeper.IsForbiddenByTokenIssuer(ctx, money, order.Sender)) {
			ordersOut = append(ordersOut, order)
		}
	}
	return ordersOut
}

func runMatch(ctx sdk.Context, midPrice sdk.Dec, ratio int, symbol string, keeper Keeper, dataHash []byte, currHeight int64) (map[string]*Order, sdk.Dec) {
	orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
	asKeeper := keeper.axk
	bxKeeper := keeper.bnk
	lowPrice := midPrice.Mul(sdk.NewDec(int64(100 - ratio))).Quo(sdk.NewDec(100))
	highPrice := midPrice.Mul(sdk.NewDec(int64(100 + ratio))).Quo(sdk.NewDec(100))

	infoForDeal := &InfoForDeal{
		bxKeeper:      bxKeeper,
		dataHash:      dataHash,
		changedOrders: make(map[string]*Order),
		context:       ctx,
		lastPrice:     sdk.NewDec(0),
	}

	// from the order book, we fetch the candidate orders for matching and filter them
	stockAndMoney := strings.Split(orderKeeper.GetSymbol(), "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	orderCandidates := orderKeeper.GetMatchingCandidates(ctx)
	orderCandidates = filterCandidates(ctx, asKeeper, orderCandidates, stock, money)

	// fill bidList and askList with wrapped orders
	bidList := make([]match.OrderForTrade, 0, len(orderCandidates))
	askList := make([]match.OrderForTrade, 0, len(orderCandidates))
	for _, orderCandidate := range orderCandidates {
		wrappedOrder := &WrappedOrder{
			order:       orderCandidate,
			infoForDeal: infoForDeal,
		}
		if wrappedOrder.order.Side == match.BID {
			bidList = append(bidList, wrappedOrder)
		} else {
			askList = append(askList, wrappedOrder)
		}
	}

	// call the match engine
	match.Match(highPrice, midPrice, lowPrice, bidList, askList)

	// both dealt orders and IOC order need further processing
	ordersForUpdate := infoForDeal.changedOrders
	for _, order := range orderKeeper.GetOrdersAtHeight(ctx, currHeight) {
		if order.TimeInForce == IOC {
			// if an IOC order is not included, we include it
			if _, ok := ordersForUpdate[order.OrderID()]; !ok {
				ordersForUpdate[order.OrderID()] = order
			}
		}
	}
	return ordersForUpdate, infoForDeal.lastPrice
}

func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	recordDay := keeper.orderClean.GetDay(ctx)
	currDay := ctx.BlockHeader().Time.Day()
	marketInfoList := keeper.GetAllMarketInfos(ctx)
	currHeight := ctx.BlockHeight()
	marketParams := keeper.GetParams(ctx)

	// if this is the first block of a new day, we clean the GTE order and there is no trade
	lifeTime := marketParams.GTEOrderLifetime
	if currDay != recordDay {
		keeper.orderClean.SetDay(ctx, currDay)
		for _, mi := range marketInfoList {
			symbol := mi.Stock + "/" + mi.Money
			orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
			removeOrderOlderThan(ctx, orderKeeper, keeper.bnk, currHeight-int64(lifeTime))
		}
		return nil
	}

	ordersForUpdateList := make([]map[string]*Order, len(marketInfoList))
	newPrices := make([]sdk.Dec, len(marketInfoList))
	for idx, mi := range marketInfoList {
		// if a token is globally forbidden, exchange it is also impossible
		if keeper.axk.IsTokenForbidden(ctx, mi.Stock) ||
			keeper.axk.IsTokenForbidden(ctx, mi.Money) {
			continue
		}
		symbol := mi.Stock + "/" + mi.Money
		dataHash := ctx.BlockHeader().DataHash
		ratio := marketParams.MaxExecutedPriceChangeRatio
		oUpdate, newPrice := runMatch(ctx, mi.LastExecutedPrice, ratio, symbol, keeper, dataHash, currHeight)
		newPrices[idx] = newPrice
		ordersForUpdateList[idx] = oUpdate
	}
	for idx, mi := range marketInfoList {
		// ignore a market if there are no orders need further processing
		if len(ordersForUpdateList[idx]) == 0 {
			continue
		}
		symbol := mi.Stock + "/" + mi.Money
		orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
		// update the order book
		for _, order := range ordersForUpdateList[idx] {
			orderKeeper.Add(ctx, order)
			if order.TimeInForce == IOC || order.LeftStock == 0 || notEnoughMoney(order) {
				removeOrder(ctx, orderKeeper, keeper.bnk, order)
			}
		}
		// if some orders dealt, update last executed price of this market
		if !newPrices[idx].IsZero() {
			mi.LastExecutedPrice = newPrices[idx]
			keeper.SetMarket(ctx, mi)
		}
	}

	// process the delist requests
	delistKeeper := NewDelistKeeper(keeper.marketKey)
	delistSymbols := delistKeeper.GetDelistSymbolsAtHeight(ctx, ctx.BlockHeight())
	for _, symbol := range delistSymbols {
		orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
		removeOrderOlderThan(ctx, orderKeeper, keeper.bnk, currHeight+1)
		keeper.RemoveMarket(ctx, symbol)
	}
	delistKeeper.RemoveDelistRequestsAtHeight(ctx, ctx.BlockHeight())

	return nil
}
