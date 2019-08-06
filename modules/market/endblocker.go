package market

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	mtype "github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/match"
)

// Some handlers which are useful when orders are matched and traded.
type InfoForDeal struct {
	bxKeeper      mtype.ExpectedBankxKeeper
	dataHash      []byte
	changedOrders map[string]*mtype.Order
	lastPrice     sdk.Dec
	context       sdk.Context
}

// returns true when a buyer's frozen money is not enough to buy LeftStock.
func notEnoughMoney(order *mtype.Order) bool {
	return order.Side == mtype.BUY &&
		order.Freeze < order.Price.Mul(sdk.NewDec(order.LeftStock)).RoundInt64()
}

// Wrapper an order with InfoForDeal, which contains useful handlers
type WrappedOrder struct {
	order       *mtype.Order
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
	if buyer.Side == mtype.SELL {
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
func unfreezeCoinsForOrder(ctx sdk.Context, bxKeeper mtype.ExpectedBankxKeeper, order *mtype.Order, feeForZeroDeal int64, feeK mtype.ExpectedChargeFeeKeeper) {
	stockAndMoney := strings.Split(order.TradingPair, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	frozenToken := stock
	if order.Side == mtype.BUY {
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
func removeOrderOlderThan(ctx sdk.Context, orderKeeper keepers.OrderKeeper, bxKeeper mtype.ExpectedBankxKeeper, feeK mtype.ExpectedChargeFeeKeeper, height int64, feeForZeroDeal int64) {
	for _, order := range orderKeeper.GetOlderThan(ctx, height) {
		removeOrder(ctx, orderKeeper, bxKeeper, feeK, order, feeForZeroDeal)
	}
}

// unfreeze the frozen token in the order and remove it from the market
func removeOrder(ctx sdk.Context, orderKeeper keepers.OrderKeeper, bxKeeper mtype.ExpectedBankxKeeper, feeK mtype.ExpectedChargeFeeKeeper, order *mtype.Order, feeForZeroDeal int64) {

	if order.Freeze != 0 || order.FrozenFee != 0 {
		unfreezeCoinsForOrder(ctx, bxKeeper, order, feeForZeroDeal, feeK)
	}
	orderKeeper.Remove(ctx, order)
}

// Iterate the candidate orders for matching, and remove the orders whose sender is forbidden by the money owner or the stock owner.
func filterCandidates(ctx sdk.Context, asKeeper mtype.ExpectedAssetStatusKeeper, ordersIn []*mtype.Order, stock, money string) []*mtype.Order {
	ordersOut := make([]*mtype.Order, 0, len(ordersIn))
	for _, order := range ordersIn {
		if !(asKeeper.IsForbiddenByTokenIssuer(ctx, stock, order.Sender) ||
			asKeeper.IsForbiddenByTokenIssuer(ctx, money, order.Sender)) {
			ordersOut = append(ordersOut, order)
		}
	}
	return ordersOut
}

func runMatch(ctx sdk.Context, midPrice sdk.Dec, ratio int, symbol string, keeper keepers.Keeper, dataHash []byte, currHeight int64) (map[string]*mtype.Order, sdk.Dec) {
	orderKeeper := keepers.NewOrderKeeper(keeper.GetMarketKey(), symbol, mtype.ModuleCdc)
	asKeeper := keeper.GetAssetKeeper()
	bxKeeper := keeper.GetBankxKeeper()
	lowPrice := midPrice.Mul(sdk.NewDec(int64(100 - ratio))).Quo(sdk.NewDec(100))
	highPrice := midPrice.Mul(sdk.NewDec(int64(100 + ratio))).Quo(sdk.NewDec(100))

	infoForDeal := &InfoForDeal{
		bxKeeper:      bxKeeper,
		dataHash:      dataHash,
		changedOrders: make(map[string]*mtype.Order),
		context:       ctx,
		lastPrice:     sdk.NewDec(0),
	}

	// from the order book, we fetch the candidate orders for matching and filter them
	stockAndMoney := strings.Split(orderKeeper.GetSymbol(), "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	orderCandidates := orderKeeper.GetMatchingCandidates(ctx)
	orderCandidates = filterCandidates(ctx, asKeeper, orderCandidates, stock, money)
	orderOldDeals := make(map[string]int64, len(orderCandidates))
	orderOldMoneys := make(map[string]int64, len(orderCandidates))
	for _, order := range orderCandidates {
		orderOldDeals[order.OrderID()] = order.DealStock
		orderOldMoneys[order.OrderID()] = order.DealMoney
	}

	// fill bidList and askList with wrapped orders
	bidList := make([]match.OrderForTrade, 0, len(orderCandidates))
	askList := make([]match.OrderForTrade, 0, len(orderCandidates))
	for _, orderCandidate := range orderCandidates {
		wrappedOrder := &WrappedOrder{
			order:       orderCandidate,
			infoForDeal: infoForDeal,
		}
		if wrappedOrder.order.Side == mtype.BID {
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
		if order.TimeInForce == mtype.IOC {
			// if an IOC order is not included, we include it
			if _, ok := ordersForUpdate[order.OrderID()]; !ok {
				ordersForUpdate[order.OrderID()] = order
			}
		}
	}

	if keeper.IsSubScribe(mtype.Topic) {
		sendFillMsg(ctx, orderOldDeals, orderOldMoneys, ordersForUpdate, ctx.BlockHeight())
	}
	return ordersForUpdate, infoForDeal.lastPrice
}

func sendFillMsg(ctx sdk.Context, orderOldDeal, orderOldMoneys map[string]int64, ordersForUpdate map[string]*mtype.Order, height int64) {
	if len(ordersForUpdate) == 0 {
		return
	}

	for id, order := range ordersForUpdate {
		oldDeal := orderOldDeal[id]
		oldMoney := orderOldMoneys[id]
		msgInfo := mtype.FillOrderInfo{
			OrderID:   id,
			Height:    height,
			Side:      order.Side,
			LeftStock: order.LeftStock,
			Freeze:    order.Freeze,
			DealStock: order.DealStock,
			DealMoney: order.DealMoney,
			CurrStock: order.DealStock - oldDeal,
			CurrMoney: order.DealMoney - oldMoney,
			Price:     order.Price.String(),
		}
		fillMsgs(ctx, mtype.FillOrderInfoKey, msgInfo)
	}
}

func filterOldOrders(oldOrders []*mtype.Order, currHeight int64, lifeTime int) []*mtype.Order {
	orders := make([]*mtype.Order, 0, len(oldOrders))
	for _, order := range oldOrders {
		if currHeight-int64(lifeTime) < int64(order.ExistBlocks) {
			continue
		}
		orders = append(orders, order)
	}

	return orders
}

func removeExpiredOrder(ctx sdk.Context, keeper keepers.Keeper, marketInfoList []mtype.MarketInfo, marketParams keepers.Params) {
	lifeTime := marketParams.GTEOrderLifetime
	currHeight := ctx.BlockHeight()
	for _, mi := range marketInfoList {
		symbol := mi.Stock + mtype.SymbolSeparator + mi.Money
		orderKeeper := keepers.NewOrderKeeper(keeper.GetMarketKey(), symbol, mtype.ModuleCdc)
		oldOrders := orderKeeper.GetOlderThan(ctx, currHeight-int64(lifeTime))
		filterOrders := filterOldOrders(oldOrders, currHeight, lifeTime)

		for _, order := range filterOrders {
			removeOrder(ctx, orderKeeper, keeper.GetBankxKeeper(), keeper, order, marketParams.FeeForZeroDeal)
			if keeper.IsSubScribe(mtype.Topic) {
				msgInfo := mtype.CancelOrderInfo{
					OrderID:        order.OrderID(),
					Side:           order.Side,
					DelReason:      mtype.CancelOrderByGteTimeOut,
					DelHeight:      ctx.BlockHeight(),
					UsedCommission: order.CalOrderFee(marketParams.FeeForZeroDeal).RoundInt64(),
					LeftStock:      order.LeftStock,
					RemainAmount:   order.Freeze,
					DealStock:      order.DealStock,
					DealMoney:      order.DealMoney,
				}
				fillMsgs(ctx, mtype.CancelOrderInfoKey, msgInfo)
			}
		}
	}
}

func removeExpiredMarket(ctx sdk.Context, keeper keepers.Keeper, marketParams keepers.Params) {
	currHeight := ctx.BlockHeight()
	currTime := ctx.BlockHeader().Time.Unix()

	// process the delist requests
	delistKeeper := keepers.NewDelistKeeper(keeper.GetMarketKey())
	delistSymbols := delistKeeper.GetDelistSymbolsBeforeTime(ctx, currTime-marketParams.MarketMinExpiredTime+1)
	for _, symbol := range delistSymbols {
		orderKeeper := keepers.NewOrderKeeper(keeper.GetMarketKey(), symbol, mtype.ModuleCdc)
		removeOrderOlderThan(ctx, orderKeeper, keeper.GetBankxKeeper(), keeper, currHeight+1, marketParams.FeeForZeroDeal)
		keeper.RemoveMarket(ctx, symbol)
	}
	delistKeeper.RemoveDelistRequestsBeforeTime(ctx, currTime-marketParams.MarketMinExpiredTime+1)
}

func EndBlocker(ctx sdk.Context, keeper keepers.Keeper) /*sdk.Tags*/ {

	marketInfoList := keeper.GetAllMarketInfos(ctx)
	currHeight := ctx.BlockHeight()
	marketParams := keeper.GetParams(ctx)

	chainID := ctx.ChainID()
	recordTime := keeper.GetOrderCleanTime(ctx)
	currTime := ctx.BlockHeader().Time.Unix()

	var needRemove bool
	if !strings.Contains(chainID, IntegrationNetSubString) {
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
		keeper.SetOrderCleanTime(ctx, currTime)
		removeExpiredOrder(ctx, keeper, marketInfoList, marketParams)
		removeExpiredMarket(ctx, keeper, marketParams)
		return //nil
	}

	ordersForUpdateList := make([]map[string]*mtype.Order, len(marketInfoList))
	newPrices := make([]sdk.Dec, len(marketInfoList))
	for idx, mi := range marketInfoList {
		// if a token is globally forbidden, exchange it is also impossible
		if keeper.IsTokenForbidden(ctx, mi.Stock) ||
			keeper.IsTokenForbidden(ctx, mi.Money) {
			continue
		}
		symbol := mi.Stock + mtype.SymbolSeparator + mi.Money
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
		orderKeeper := keepers.NewOrderKeeper(keeper.GetMarketKey(), symbol, mtype.ModuleCdc)
		// update the order book
		for _, order := range ordersForUpdateList[idx] {
			orderKeeper.Add(ctx, order)
			if order.TimeInForce == mtype.IOC || order.LeftStock == 0 || notEnoughMoney(order) {
				if keeper.IsSubScribe(mtype.Topic) {
					sendOrderMsg(ctx, order, ctx.BlockHeight(), marketParams.FeeForZeroDeal, keeper)
				}
				removeOrder(ctx, orderKeeper, keeper.GetBankxKeeper(), keeper, order, marketParams.FeeForZeroDeal)
			}
		}
		// if some orders dealt, update last executed price of this market
		if !newPrices[idx].IsZero() {
			mi.LastExecutedPrice = newPrices[idx]
			keeper.SetMarket(ctx, mi)
		}
	}

}

func sendOrderMsg(ctx sdk.Context, order *mtype.Order, height int64, feeForZeroDeal int64, keeper keepers.Keeper) {
	msgInfo := mtype.CancelOrderInfo{
		OrderID:        order.OrderID(),
		Side:           order.Side,
		DelHeight:      height,
		UsedCommission: order.CalOrderFee(feeForZeroDeal).RoundInt64(),
		LeftStock:      order.LeftStock,
		RemainAmount:   order.Freeze,
		DealStock:      order.DealStock,
		DealMoney:      order.DealMoney,
	}
	if order.TimeInForce == mtype.IOC {
		msgInfo.DelReason = mtype.CancelOrderByIocType
	} else if order.LeftStock == 0 {
		msgInfo.DelReason = mtype.CancelOrderByAllFilled
	} else if notEnoughMoney(order) {
		msgInfo.DelReason = mtype.CancelOrderByNoEnoughMoney
	} else {
		msgInfo.DelReason = mtype.CancelOrderByNotKnow
	}

	fillMsgs(ctx, mtype.CancelOrderInfoKey, msgInfo)
}
