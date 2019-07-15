package market

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/coinexchain/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/match"
)

const (
	testNetSubString = "coinexdex-test"
	mainNetSubString = "coinexdex-main"
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
	stockAndMoney := strings.Split(buyer.TradingPair, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	// buyer and seller will exchange stockCoins and moneyCoins
	stockCoins := sdk.Coins{sdk.NewCoin(stock, sdk.NewInt(amount))}
	moneyAmount := price.MulInt(sdk.NewInt(amount)).TruncateInt64()
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
func unfreezeCoinsForOrder(ctx sdk.Context, bxKeeper ExpectedBankxKeeper, order *Order, feeForZeroDeal int64, feeK ExpectedChargeFeeKeeper) {
	stockAndMoney := strings.Split(order.TradingPair, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	frozenToken := stock
	if order.Side == match.BUY {
		frozenToken = money
	}

	coins := sdk.Coins([]sdk.Coin{sdk.NewCoin(frozenToken, sdk.NewInt(order.Freeze))})
	bxKeeper.UnFreezeCoins(ctx, order.Sender, coins)

	if order.FrozenFee != 0 {
		coins = sdk.Coins([]sdk.Coin{sdk.NewCoin(types.CET, sdk.NewInt(order.FrozenFee))})
		bxKeeper.UnFreezeCoins(ctx, order.Sender, coins)
		actualFee := order.CalOrderFee(feeForZeroDeal)
		if err := feeK.SubtractFeeAndCollectFee(ctx, order.Sender, types.NewCetCoins(actualFee.RoundInt64())); err != nil {
			panic(err)
		}
	}
}

// remove the orders whose age are older than height
func removeOrderOlderThan(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, feeK ExpectedChargeFeeKeeper, height int64, feeForZeroDeal int64) {
	for _, order := range orderKeeper.GetOlderThan(ctx, height) {
		removeOrder(ctx, orderKeeper, bxKeeper, feeK, order, feeForZeroDeal)
	}
}

// unfreeze the frozen token in the order and remove it from the market
func removeOrder(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, feeK ExpectedChargeFeeKeeper, order *Order, feeForZeroDeal int64) {

	if order.Freeze != 0 || order.FrozenFee != 0 {
		unfreezeCoinsForOrder(ctx, bxKeeper, order, feeForZeroDeal, feeK)
	}
	orderKeeper.Remove(ctx, order)
}

// Iterate the candidate orders for matching, and remove the orders whose sender is forbidden by the money owner or the stock owner.
func filterCandidates(ctx sdk.Context, asKeeper ExpectedAssetStatusKeeper, ordersIn []*Order, stock, money string) []*Order {
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
	orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, ModuleCdc)
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
	orderOldDeals := make(map[string]int64, len(orderCandidates))
	for _, order := range orderCandidates {
		orderOldDeals[order.OrderID()] = order.DealStock
	}

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

	sendFillMsg(orderOldDeals, ordersForUpdate, ctx.BlockHeight(), infoForDeal.lastPrice, keeper)
	return ordersForUpdate, infoForDeal.lastPrice
}

func sendFillMsg(orderOldDeal map[string]int64, ordersForUpdate map[string]*Order, height int64, price sdk.Dec, keeper Keeper) {
	if len(ordersForUpdate) == 0 {
		return
	}

	for id, order := range ordersForUpdate {
		oldDeal := orderOldDeal[id]
		msgInfo := FillOrderInfo{
			OrderID:   id,
			Height:    height,
			LeftStock: order.LeftStock,
			Freeze:    order.Freeze,
			DealStock: order.DealStock,
			DealMoney: order.DealMoney,
			CurrStock: order.DealStock - oldDeal,
			Price:     price.String(),
		}
		keeper.SendMsg(FillOrderInfoKey, msgInfo)
	}
}

func filterOldOrders(oldOrders []*Order, currHeight int64, lifeTime int) []*Order {
	orders := make([]*Order, 0, len(oldOrders))
	for _, order := range oldOrders {
		if currHeight-int64(lifeTime) < int64(order.ExistBlocks) {
			continue
		}
		orders = append(orders, order)
	}

	return orders
}

func removeExpiredOrder(ctx sdk.Context, keeper Keeper, marketInfoList []MarketInfo, marketParams Params) {
	lifeTime := marketParams.GTEOrderLifetime
	currHeight := ctx.BlockHeight()
	for _, mi := range marketInfoList {
		symbol := mi.Stock + SymbolSeparator + mi.Money
		orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, ModuleCdc)
		oldOrders := orderKeeper.GetOlderThan(ctx, currHeight-int64(lifeTime))
		filterOrders := filterOldOrders(oldOrders, currHeight, lifeTime)

		for _, order := range filterOrders {
			removeOrder(ctx, orderKeeper, keeper.bnk, keeper, order, marketParams.FeeForZeroDeal)

			msgInfo := CancelOrderInfo{
				OrderID:        order.OrderID(),
				DelReason:      CancelOrderByGteTimeOut,
				DelHeight:      ctx.BlockHeight(),
				UsedCommission: order.CalOrderFee(marketParams.FeeForZeroDeal).RoundInt64(),
				LeftStock:      order.LeftStock,
				RemainAmount:   order.Freeze,
				DealStock:      order.DealStock,
				DealMoney:      order.DealMoney,
			}
			keeper.SendMsg(CancelOrderInfoKey, msgInfo)
		}
	}
}

func removeExpiredMarket(ctx sdk.Context, keeper Keeper, marketParams Params) {
	currHeight := ctx.BlockHeight()
	currTime := ctx.BlockHeader().Time.Unix()

	// process the delist requests
	delistKeeper := NewDelistKeeper(keeper.marketKey)
	delistSymbols := delistKeeper.GetDelistSymbolsBeforeTime(ctx, currTime-marketParams.MarketMinExpiredTime+1)
	for _, symbol := range delistSymbols {
		orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, ModuleCdc)
		removeOrderOlderThan(ctx, orderKeeper, keeper.bnk, keeper, currHeight+1, marketParams.FeeForZeroDeal)
		keeper.RemoveMarket(ctx, symbol)
	}
	delistKeeper.RemoveDelistRequestsBeforeTime(ctx, currTime-marketParams.MarketMinExpiredTime+1)
}

func EndBlocker(ctx sdk.Context, keeper Keeper) /*sdk.Tags*/ {

	marketInfoList := keeper.GetAllMarketInfos(ctx)
	currHeight := ctx.BlockHeight()
	marketParams := keeper.GetParams(ctx)

	chainID := ctx.ChainID()
	recordTime := keeper.orderClean.GetUnixTime(ctx)
	currTime := ctx.BlockHeader().Time.Unix()

	var needRemove bool
	if strings.Contains(chainID, mainNetSubString) || strings.Contains(chainID, testNetSubString) {
		if time.Unix(recordTime, 0).Day() != time.Unix(currTime, 0).Day() {
			needRemove = true
		}
	} else {
		if time.Unix(recordTime, 0).Minute() != time.Unix(currTime, 0).Minute() {
			needRemove = true
		}
	}

	// if this is the first block of a new day, we clean the GTE order and there is no trade
	if needRemove {
		keeper.orderClean.SetUnixTime(ctx, currTime)
		removeExpiredOrder(ctx, keeper, marketInfoList, marketParams)
		removeExpiredMarket(ctx, keeper, marketParams)
		return //nil
	}

	ordersForUpdateList := make([]map[string]*Order, len(marketInfoList))
	newPrices := make([]sdk.Dec, len(marketInfoList))
	for idx, mi := range marketInfoList {
		// if a token is globally forbidden, exchange it is also impossible
		if keeper.axk.IsTokenForbidden(ctx, mi.Stock) ||
			keeper.axk.IsTokenForbidden(ctx, mi.Money) {
			continue
		}
		symbol := mi.Stock + SymbolSeparator + mi.Money
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
		orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, ModuleCdc)
		// update the order book
		for _, order := range ordersForUpdateList[idx] {
			orderKeeper.Add(ctx, order)
			if order.TimeInForce == IOC || order.LeftStock == 0 || notEnoughMoney(order) {
				sendOrderMsg(order, ctx.BlockHeight(), marketParams.FeeForZeroDeal, keeper)
				removeOrder(ctx, orderKeeper, keeper.bnk, keeper, order, marketParams.FeeForZeroDeal)
			}
		}
		// if some orders dealt, update last executed price of this market
		if !newPrices[idx].IsZero() {
			mi.LastExecutedPrice = newPrices[idx]
			keeper.SetMarket(ctx, mi)
		}
	}

	//return nil
}

func sendOrderMsg(order *Order, height int64, feeForZeroDeal int64, keeper Keeper) {
	msgInfo := CancelOrderInfo{
		OrderID:        order.OrderID(),
		DelHeight:      height,
		UsedCommission: order.CalOrderFee(feeForZeroDeal).RoundInt64(),
		LeftStock:      order.LeftStock,
		RemainAmount:   order.Freeze,
		DealStock:      order.DealStock,
		DealMoney:      order.DealMoney,
	}
	if order.TimeInForce == IOC {
		msgInfo.DelReason = CancelOrderByIocType
	} else if order.LeftStock == 0 {
		msgInfo.DelReason = CancelOrderByAllFilled
	} else if notEnoughMoney(order) {
		msgInfo.DelReason = CancelOrderByNoEnoughMoney
	} else {
		msgInfo.DelReason = CancelOrderByNotKnow
	}

	keeper.SendMsg(CancelOrderInfoKey, msgInfo)
}
