package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
	simulationx "github.com/coinexchain/dex/simulation"
	dex "github.com/coinexchain/dex/types"
)

// TODO

func SimulateMsgCreateTradingPair(k keepers.Keeper, ask asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		fromAcc := simulation.RandomAcc(r, accs)

		msg, err := createMsgCreateTradingPair(r, ctx, k, ask, fromAcc.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("error create msg")
		}

		handler := market.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyCreateTradingPair(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("trading pair creation failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func createMsgCreateTradingPair(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ask asset.Keeper, fromAddr sdk.AccAddress) (types.MsgCreateTradingPair, error) {

	stock, money := randomTradingPair(r, ctx, ask)
	if stock == "" || money == "" {
		return types.MsgCreateTradingPair{}, nil
	}

	precision := r.Intn(types.MaxTokenPricePrecision + 1)

	if simulationx.RandomBool(r) {
		fromAddr = ask.GetToken(ctx, stock).GetOwner()
	}

	if money != dex.CET && stock != dex.CET {
		if _, err := k.GetMarketInfo(ctx, types.GetSymbol(stock, dex.CET)); err != nil {
			money = dex.CET
		}
	}

	msg := types.MsgCreateTradingPair{
		Stock:          stock,
		Money:          money,
		Creator:        fromAddr,
		PricePrecision: byte(precision),
	}
	if msg.ValidateBasic() != nil {
		return types.MsgCreateTradingPair{}, fmt.Errorf("msg expected to pass validation check")
	}

	return msg, nil
}
func randomTradingPair(r *rand.Rand, ctx sdk.Context, ask asset.Keeper) (stock, money string) {
	tokenList := ask.GetAllTokens(ctx)
	if len(tokenList) < 2 {
		return
	}
	stock = tokenList[simulationx.GetRandomElemIndex(r, len(tokenList))].GetSymbol()
	for {
		money = tokenList[simulationx.GetRandomElemIndex(r, len(tokenList))].GetSymbol()
		if stock != money {
			break
		}
	}
	return
}

func verifyCreateTradingPair(ctx sdk.Context, k keepers.Keeper, msg types.MsgCreateTradingPair) bool {
	tradingPair, err := k.GetMarketInfo(ctx, msg.GetSymbol())
	return err == nil &&
		tradingPair.PricePrecision == msg.PricePrecision
}

func SimulateMsgCancelTradingPair(k keepers.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {
		fromAcc := simulation.RandomAcc(r, accs)

		msg, err := createMsgCancelTradingPair(r, ctx, k, fromAcc.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		handler := market.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}
		if time.Now().Unix() < msg.EffectiveTime {
			return simulation.NewOperationMsg(msg, ok, ""), []simulation.FutureOperation{
				{
					BlockTime: time.Unix(msg.EffectiveTime, 0),
					Op:        SimulateVerifyCancelTradingPair(k, msg),
				},
			}, nil
		}
		ok = verifyCancelTradingPair(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("trading pair cancel failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func createMsgCancelTradingPair(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, fromAddr sdk.AccAddress) (types.MsgCancelTradingPair, error) {

	tradingPair, err := randomExistedTradingPair(r, ctx, k)
	if err != nil {
		return types.MsgCancelTradingPair{}, fmt.Errorf("no trading pair to cancel")
	}
	timeStamp := simulation.RandTimestamp(r)
	msg := types.MsgCancelTradingPair{
		Sender:        fromAddr,
		TradingPair:   tradingPair.GetSymbol(),
		EffectiveTime: timeStamp.Unix(),
	}
	if msg.ValidateBasic() != nil {
		return types.MsgCancelTradingPair{}, fmt.Errorf("msg expected to pass validation check")
	}
	return msg, nil
}
func verifyCancelTradingPair(ctx sdk.Context, k keepers.Keeper, msg types.MsgCancelTradingPair) bool {
	_, err := k.GetMarketInfo(ctx, msg.TradingPair)
	return err != nil
}

func SimulateVerifyCancelTradingPair(k keepers.Keeper, msg types.MsgCancelTradingPair) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {
		ok := verifyCancelTradingPair(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("trading pair cancel failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func SimulateMsgModifyPricePrecision(k keepers.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {
		fromAcc := simulation.RandomAcc(r, accs)

		msg, err := createMsgModifyPricePrecision(r, ctx, k, fromAcc.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		handler := market.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyModifyPricePrecision(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("trading pair price precision modification failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func createMsgModifyPricePrecision(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, fromAddr sdk.AccAddress) (types.MsgModifyPricePrecision, error) {

	tradingPair, err := randomExistedTradingPair(r, ctx, k)
	if err != nil {
		return types.MsgModifyPricePrecision{}, fmt.Errorf("no trading pair to modify price precision")
	}
	newPrecision := r.Intn(types.MaxTokenPricePrecision + 1)
	msg := types.MsgModifyPricePrecision{
		TradingPair:    tradingPair.GetSymbol(),
		PricePrecision: byte(newPrecision),
		Sender:         fromAddr,
	}
	if msg.ValidateBasic() != nil {
		return types.MsgModifyPricePrecision{}, fmt.Errorf("msg expected to pass validation check")
	}
	return msg, nil

}
func verifyModifyPricePrecision(ctx sdk.Context, k keepers.Keeper, msg types.MsgModifyPricePrecision) bool {
	tradingPair, err := k.GetMarketInfo(ctx, msg.TradingPair)
	return err == nil && tradingPair.PricePrecision == msg.PricePrecision
}

func SimulateMsgCreateOrder(k keepers.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		msg, err := createMsgCreateOrder(r, ctx, k, ak, accs)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		handler := market.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyCreateOrder(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("order creation failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func createMsgCreateOrder(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ak auth.AccountKeeper, accs []simulation.Account) (msg types.MsgCreateOrder, err error) {

	var tradingPair types.MarketInfo
	var denom string
	var side int
	if simulationx.RandomBool(r) {
		side = types.BUY
	} else {
		side = types.SELL
	}

	//randomly generate trading pair
	tradingPair, err = randomExistedTradingPair(r, ctx, k)
	if side == types.BUY {
		denom = tradingPair.Money
	} else {
		denom = tradingPair.Stock
	}
	if err != nil {
		err = fmt.Errorf("no existed trading pair to create order")
		return
	}

	//generate from account
	var fromCoins sdk.Coins
	var fromAddr sdk.AccAddress
	var fromAcc auth.Account
	for {
		fromAddr = simulation.RandomAcc(r, accs).Address
		fromAcc = ak.GetAccount(ctx, fromAddr)
		fromCoins = fromAcc.GetCoins()
		if !fromCoins.AmountOf(denom).IsZero() {
			break
		}
	}

	//randomly generate price & quantity to trade
	coinsHold := fromCoins.AmountOf(denom)
	var quantity, price int64
	var precision byte
	for {
		precision = byte(r.Intn(int(tradingPair.PricePrecision) + 1))
		quantity = r.Int63n(coinsHold.Int64())
		price = r.Int63n(int64(math.Pow10(int(precision))))
		if side == types.SELL || calculateAmount(price, quantity, precision).LT(coinsHold.ToDec()) {
			break
		}
	}

	var timeInforce int
	if simulationx.RandomBool(r) {
		timeInforce = types.IOC
	} else {
		timeInforce = types.GTE
	}

	msg = types.MsgCreateOrder{
		Sender:         fromAddr,
		Identify:       byte(r.Intn(255 + 1)),
		TradingPair:    tradingPair.GetSymbol(),
		OrderType:      types.LimitOrder,
		PricePrecision: precision,
		Price:          price,
		Quantity:       quantity,
		Side:           byte(side),
		TimeInForce:    timeInforce,
	}
	if msg.ValidateBasic() != nil {
		return types.MsgCreateOrder{}, fmt.Errorf("msg expected to pass validation check")
	}
	return msg, nil

}
func calculateAmount(price, quantity int64, pricePrecision byte) sdk.Dec {
	actualPrice := sdk.NewDec(price).Quo(sdk.NewDec(int64(math.Pow10(int(pricePrecision)))))
	money := actualPrice.Mul(sdk.NewDec(quantity))
	return money.Add(sdk.NewDec(types.ExtraFrozenMoney)).Ceil()
}
func randomExistedTradingPair(r *rand.Rand, ctx sdk.Context, k keepers.Keeper) (tradingPair types.MarketInfo, err error) {
	tradingPairs := k.GetAllMarketInfos(ctx)
	if len(tradingPairs) == 0 {
		err = fmt.Errorf("no existed trading pair")
		return
	}
	tradingPair = tradingPairs[r.Intn(len(tradingPairs))]
	return
}
func verifyCreateOrder(ctx sdk.Context, k keepers.Keeper, msg types.MsgCreateOrder) bool {

	ork := keepers.NewGlobalOrderKeeper(k.GetMarketKey(), types.ModuleCdc)
	orderID, err := types.AssemblyOrderID(msg.Sender.String(), 0, msg.Identify)
	if err != nil {
		return false
	}
	order := ork.QueryOrder(ctx, orderID)
	return order.Sender.Equals(msg.Sender) &&
		order.TimeInForce == msg.TimeInForce &&
		order.Side == msg.Side &&
		order.Quantity == msg.Quantity &&
		order.Price.Equal(sdk.NewDec(msg.Price).Quo(sdk.NewDec(int64(math.Pow10(int(msg.PricePrecision)))))) &&
		order.TradingPair == msg.TradingPair &&
		order.OrderType == msg.OrderType
}

func SimulateMsgCancelOrder(k keepers.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		msg, err := createMsgCancelOrder(r, ctx, k)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		handler := market.NewHandler(k)
		ok := simulationx.SimulateHandleMsg(msg, handler, ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyCancelOrder(ctx, k, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("order cancel failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func createMsgCancelOrder(r *rand.Rand, ctx sdk.Context, k keepers.Keeper) (types.MsgCancelOrder, error) {
	orders := k.GetAllOrders(ctx)
	if len(orders) == 0 {
		return types.MsgCancelOrder{}, fmt.Errorf("no order to cancel")
	}
	orderToCancel := orders[r.Intn(len(orders))]
	msg := types.MsgCancelOrder{
		Sender:  orderToCancel.Sender,
		OrderID: orderToCancel.OrderID(),
	}
	if msg.ValidateBasic() != nil {
		return types.MsgCancelOrder{}, fmt.Errorf("msg expected to pass validation check")
	}
	return msg, nil
}

func verifyCancelOrder(ctx sdk.Context, k keepers.Keeper, msg types.MsgCancelOrder) bool {

	ork := keepers.NewGlobalOrderKeeper(k.GetMarketKey(), types.ModuleCdc)
	order := ork.QueryOrder(ctx, msg.OrderID)
	return order == nil
}
