package market

import (
	"crypto/sha256"
	"github.com/coinexchain/dex/modules/market/match"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type InfoForDeal struct {
	bxKeeper        ExpectedBankxKeeper
	dataHash        []byte
	fulfilledOrders map[string]*Order
	lastPrice       sdk.Dec
	context         sdk.Context
}

type WrappedOrder struct {
	order       *Order
	infoForDeal *InfoForDeal
}

func (wo *WrappedOrder) GetPrice() sdk.Dec {
	return wo.order.Price
}

func (wo *WrappedOrder) GetAmount() int64 {
	return wo.order.Quantity
}

func (wo *WrappedOrder) GetHeight() int64 {
	return wo.order.Height
}

func (wo *WrappedOrder) GetType() int {
	return int(wo.order.OrderType)
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
	if buyer.OrderType == match.SELL {
		buyer, seller = other.order, wo.order
	}
	stockAndMoney := strings.Split(buyer.Symbol, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	stockCoins := sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(amount)))
	moneyAmount := price.MulInt(sdk.NewInt(amount)).RoundInt64()
	moneyCoins := sdk.NewCoins(sdk.NewCoin(money, sdk.NewInt(moneyAmount)))

	buyer.LeftStock -= amount
	seller.LeftStock -= amount
	buyer.Freeze -= moneyAmount
	seller.Freeze -= amount
	buyer.DealStock += amount
	seller.DealStock += amount
	buyer.DealMoney += moneyAmount
	seller.DealMoney += moneyAmount
	ctx := wo.infoForDeal.context
	wo.infoForDeal.bxKeeper.UnFreezeCoins(ctx, seller.Sender, stockCoins)
	wo.infoForDeal.bxKeeper.SendCoins(ctx, seller.Sender, buyer.Sender, stockCoins)
	wo.infoForDeal.bxKeeper.UnFreezeCoins(ctx, buyer.Sender, moneyCoins)
	wo.infoForDeal.bxKeeper.SendCoins(ctx, buyer.Sender, seller.Sender, moneyCoins)
	if buyer.LeftStock == 0 {
		wo.infoForDeal.fulfilledOrders[buyer.OrderID()] = buyer
	}
	if seller.LeftStock == 0 {
		wo.infoForDeal.fulfilledOrders[seller.OrderID()] = seller
	}
	wo.infoForDeal.lastPrice = price
}

func unfreezeCoinsForOrder(ctx sdk.Context, bxKeeper ExpectedBankxKeeper, order *Order) {
	stockAndMoney := strings.Split(order.Symbol, "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	frozenToken := stock
	if order.OrderType == match.BUY {
		frozenToken = money
	}
	coins := sdk.NewCoins(sdk.NewCoin(frozenToken, sdk.NewInt(order.Freeze)))
	bxKeeper.UnFreezeCoins(ctx, order.Sender, coins)
}

func removeOrderOlderThan(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, height int64) {
	for _, order := range orderKeeper.GetOlderThan(ctx, height) {
		removeOrder(ctx, orderKeeper, bxKeeper, order)
	}
}

func removeOrder(ctx sdk.Context, orderKeeper OrderKeeper, bxKeeper ExpectedBankxKeeper, order *Order) {
	if order.Freeze != 0 {
		unfreezeCoinsForOrder(ctx, bxKeeper, order)
	}
	orderKeeper.Remove(ctx, order)
}

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

func matchForMarket(ctx sdk.Context, midPrice sdk.Dec, symbol string, keeper Keeper, dataHash []byte, currHeight int64) (toBeDeleted map[string]*Order, lastPrice sdk.Dec) {
	orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
	asKeeper := keeper.axk
	bxKeeper := keeper.bnk
	lowPrice := midPrice.Mul(sdk.NewDec(100 - MaxExecutedPriceChangeRatio)).Quo(sdk.NewDec(100))
	highPrice := midPrice.Mul(sdk.NewDec(100 + MaxExecutedPriceChangeRatio)).Quo(sdk.NewDec(100))

	infoForDeal := &InfoForDeal{
		bxKeeper:        bxKeeper,
		dataHash:        dataHash,
		fulfilledOrders: make(map[string]*Order),
		context:         ctx,
		lastPrice:       sdk.NewDec(0),
	}
	stockAndMoney := strings.Split(orderKeeper.GetSymbol(), "/")
	stock, money := stockAndMoney[0], stockAndMoney[1]
	orderCandidates := orderKeeper.GetMatchingCandidates(ctx)
	orderCandidates = filterCandidates(ctx, asKeeper, orderCandidates, stock, money)
	bidList := make([]match.OrderForTrade, 0, len(orderCandidates))
	askList := make([]match.OrderForTrade, 0, len(orderCandidates))
	for _, orderCandidate := range orderCandidates {
		wrappedOrder := &WrappedOrder{
			order:       orderCandidate,
			infoForDeal: infoForDeal,
		}
		if wrappedOrder.order.OrderType == match.BID {
			bidList = append(bidList, wrappedOrder)
		} else {
			askList = append(askList, wrappedOrder)
		}
	}
	match.Match(highPrice, midPrice, lowPrice, bidList, askList)
	tobeDel := infoForDeal.fulfilledOrders
	for _, order := range orderKeeper.GetOrdersAtHeight(ctx, currHeight) {
		if order.OrderType == IOC {
			tobeDel[order.OrderID()] = order
		}
	}
	return tobeDel, infoForDeal.lastPrice
}

func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	recordDay := keeper.orderClean.GetDay(ctx)
	currDay := ctx.BlockHeader().Time.Day()
	marketInfoList := keeper.GetAllMarketInfos(ctx)
	toBeDeletedOrderMaps := make([]map[string]*Order, len(marketInfoList))
	newPrices := make([]sdk.Dec, len(marketInfoList))
	currHeight := ctx.BlockHeight()
	if currDay != recordDay {
		keeper.orderClean.SetDay(ctx, currDay)
		for _, mi := range marketInfoList {
			symbol := mi.Stock + "/" + mi.Money
			orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
			removeOrderOlderThan(ctx, orderKeeper, keeper.bnk, currHeight-GTEOrderLifetime)
		}
	} else {
		for idx, mi := range marketInfoList {
			if keeper.axk.IsTokenFrozen(ctx, mi.Stock) ||
				keeper.axk.IsTokenFrozen(ctx, mi.Money) {
				continue
			}
			symbol := mi.Stock + "/" + mi.Money
			dataHash := ctx.BlockHeader().DataHash
			toDel, newPrice := matchForMarket(ctx, mi.LastExecutedPrice, symbol, keeper, dataHash, currHeight)
			if !newPrice.IsZero() {
				newPrices[idx] = newPrice
				toBeDeletedOrderMaps[idx] = toDel
			}
		}
		for idx, mi := range marketInfoList {
			if toBeDeletedOrderMaps[idx] == nil {
				continue
			}
			symbol := mi.Stock + "/" + mi.Money
			orderKeeper := NewOrderKeeper(keeper.marketKey, symbol, msgCdc)
			for _, order := range toBeDeletedOrderMaps[idx] {
				removeOrder(ctx, orderKeeper, keeper.bnk, order)
			}
			mi.LastExecutedPrice = newPrices[idx]
			keeper.SetMarket(ctx, mi)
		}
	}
	return nil
}
