package market

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k keepers.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateTradingPair:
			return handleMsgCreateTradingPair(ctx, msg, k)
		case types.MsgCreateOrder:
			return handleMsgCreateOrder(ctx, msg, k)
		case types.MsgCancelOrder:
			return handleMsgCancelOrder(ctx, msg, k)
		case types.MsgCancelTradingPair:
			return handleMsgCancelTradingPair(ctx, msg, k)
		case types.MsgModifyPricePrecision:
			return handleMsgModifyPricePrecision(ctx, msg, k)
		default:
			return dex.ErrUnknownRequest(ModuleName, msg)
		}
	}
}

func handleMsgCreateTradingPair(ctx sdk.Context, msg types.MsgCreateTradingPair, keeper keepers.Keeper) sdk.Result {
	if err := checkMsgCreateTradingPair(ctx, msg, keeper); err != nil {
		return err.Result()
	}

	var orderPrecision byte
	if msg.OrderPrecision <= types.MaxOrderPrecision {
		orderPrecision = msg.OrderPrecision
	}
	info := types.MarketInfo{
		Stock:             msg.Stock,
		Money:             msg.Money,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: sdk.ZeroDec(),
		OrderPrecision:    orderPrecision,
	}

	if err := keeper.SetMarket(ctx, info); err != nil {
		// only MarshalBinaryBare can cause error here, which is impossible in production
		return err.Result()
	}

	param := keeper.GetParams(ctx)
	if err := keeper.SubtractFeeAndCollectFee(ctx, msg.Creator, param.CreateMarketFee); err != nil {
		// CreateMarketFee has been checked with HasCoins in checkMsgCreateTradingPair
		// this clause will not execute in production
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeMarket,
			sdk.NewAttribute(AttributeKeyTradingPair, msg.GetSymbol()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeKeyStock, msg.Stock),
			sdk.NewAttribute(AttributeKeyMoney, msg.Money),
			sdk.NewAttribute(AttributeKeySender, msg.Creator.String()),
			sdk.NewAttribute(AttributeKeyPricePrecision, strconv.Itoa(int(info.PricePrecision))),
			sdk.NewAttribute(AttributeKeyLastExecutePrice, info.LastExecutedPrice.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper keepers.Keeper, key string, msg interface{}) {
	if keeper.IsSubScribed(types.Topic) {
		msgqueue.FillMsgs(ctx, key, msg)
	}
}

func checkMsgCreateTradingPair(ctx sdk.Context, msg types.MsgCreateTradingPair, keeper keepers.Keeper) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	if _, err := keeper.GetMarketInfo(ctx, msg.GetSymbol()); err == nil {
		return types.ErrRepeatTradingPair()
	}

	if !keeper.IsTokenExists(ctx, msg.Money) || !keeper.IsTokenExists(ctx, msg.Stock) {
		return types.ErrTokenNoExist()
	}

	if !keeper.IsTokenIssuer(ctx, msg.Stock, msg.Creator) {
		return types.ErrInvalidTokenIssuer()
	}

	if msg.Money != dex.CET && msg.Stock != dex.CET {
		if _, err := keeper.GetMarketInfo(ctx, GetSymbol(msg.Stock, dex.CET)); err != nil {
			return types.ErrNotListedAgainstCet(msg.Stock)
		}
	}

	marketParams := keeper.GetParams(ctx)
	if !keeper.HasCoins(ctx, msg.Creator, dex.NewCetCoins(marketParams.CreateMarketFee)) {
		return types.ErrInsufficientCoins()
	}

	return nil
}

func calFrozenFeeInOrder(ctx sdk.Context, marketParams types.Params, keeper keepers.Keeper, msg types.MsgCreateOrder) (int64, sdk.Error) {
	var frozenFeeDec sdk.Dec
	stock, _ := SplitSymbol(msg.TradingPair)

	// Calculate the fee when stock is cet
	rate := sdk.NewDec(marketParams.MarketFeeRate)
	div := sdk.NewDec(int64(math.Pow10(types.MarketFeeRatePrecision)))
	if stock == dex.CET {
		frozenFeeDec = sdk.NewDec(msg.Quantity).Mul(rate).Quo(div).Ceil()
	} else {
		stockSepCet := GetSymbol(stock, dex.CET)
		marketInfo, err := keeper.GetMarketInfo(ctx, stockSepCet)
		if err != nil || marketInfo.LastExecutedPrice.IsZero() {
			frozenFeeDec = sdk.NewDec(marketParams.FixedTradeFee)
		} else {
			totalPriceInCet := marketInfo.LastExecutedPrice.Mul(sdk.NewDec(msg.Quantity))
			frozenFeeDec = totalPriceInCet.Mul(rate).Quo(div).Ceil()
		}
	}
	if frozenFeeDec.GT(sdk.NewDec(types.MaxOrderAmount)) {
		return 0, types.ErrInvalidOrderAmount("The frozen fee is too large")
	}
	frozenFee := frozenFeeDec.RoundInt64()
	if frozenFee < marketParams.MarketFeeMin {
		return 0, types.ErrInvalidOrderCommission(fmt.Sprintf("%d", frozenFee))
	}

	return frozenFee, nil
}

func calFeatureFeeForExistBlocks(msg types.MsgCreateOrder, marketParam types.Params) int64 {
	if msg.TimeInForce == types.IOC {
		return 0
	}
	if msg.ExistBlocks <= marketParam.GTEOrderLifetime {
		return 0
	}
	quotient := (msg.ExistBlocks + marketParam.GTEOrderLifetime - 1) / marketParam.GTEOrderLifetime
	return int64(quotient-1) * marketParam.GTEOrderFeatureFeeByBlocks
}

func handleFeeForCreateOrder(ctx sdk.Context, keeper keepers.Keeper, amount int64, denom string,
	sender sdk.AccAddress, frozenFee, featureFee int64) sdk.Error {
	coin := sdk.NewCoin(denom, sdk.NewInt(amount))
	if err := keeper.FreezeCoins(ctx, sender, sdk.Coins{coin}); err != nil {
		return err
	}
	if frozenFee != 0 {
		frozenFeeAsCet := sdk.Coins{sdk.NewCoin(dex.CET, sdk.NewInt(frozenFee))}
		if err := keeper.FreezeCoins(ctx, sender, frozenFeeAsCet); err != nil {
			return err
		}
	}
	if featureFee != 0 {
		if err := keeper.SubtractFeeAndCollectFee(ctx, sender, featureFee); err != nil {
			return err
		}
	}
	return nil
}

func sendCreateOrderMsg(ctx sdk.Context, keeper keepers.Keeper, order types.Order, featureFee int64) {
	// send msg to kafka
	msgInfo := types.CreateOrderInfo{
		OrderID:     order.OrderID(),
		Sender:      order.Sender.String(),
		TradingPair: order.TradingPair,
		OrderType:   order.OrderType,
		Price:       order.Price,
		Quantity:    order.Quantity,
		Side:        order.Side,
		TimeInForce: order.TimeInForce,
		Height:      order.Height,
		FrozenFee:   order.FrozenFee,
		Freeze:      order.Freeze,
		FeatureFee:  featureFee,
	}
	fillMsgQueue(ctx, keeper, types.CreateOrderInfoKey, msgInfo)
}

func handleMsgCreateOrder(ctx sdk.Context, msg types.MsgCreateOrder, keeper keepers.Keeper) sdk.Result {
	stock, money := SplitSymbol(msg.TradingPair)
	denom := stock
	amount := msg.Quantity
	if msg.Side == types.BUY {
		denom = money
		tmpAmount, err := calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision)
		if err != nil {
			return types.ErrInvalidOrderAmount("The frozen fee is too large").Result()
		}
		amount = tmpAmount.RoundInt64()
	}
	if amount > types.MaxOrderAmount {
		return types.ErrInvalidOrderAmount("The frozen fee is too large").Result()
	}

	seq, err := keeper.QuerySeqWithAddr(ctx, msg.Sender)
	if err != nil {
		return err.Result()
	}
	marketParams := keeper.GetParams(ctx)
	frozenFee, err := calFrozenFeeInOrder(ctx, marketParams, keeper, msg)
	if err != nil {
		return err.Result()
	}

	featureFee := calFeatureFeeForExistBlocks(msg, marketParams)
	totalFee := frozenFee + featureFee
	if featureFee > types.MaxOrderAmount ||
		frozenFee > types.MaxOrderAmount ||
		totalFee > types.MaxOrderAmount {
		return types.ErrInvalidOrderAmount("The frozen fee is too large").Result()
	}
	if ret := checkMsgCreateOrder(ctx, keeper, msg, totalFee, amount, denom, seq); !ret.IsOK() {
		return ret
	}
	existBlocks := msg.ExistBlocks
	if existBlocks == 0 && msg.TimeInForce == GTE {
		existBlocks = marketParams.GTEOrderLifetime
	}

	order := types.Order{
		Sender:      msg.Sender,
		Sequence:    seq,
		Identify:    msg.Identify,
		TradingPair: msg.TradingPair,
		OrderType:   msg.OrderType,
		Price:       sdk.NewDec(msg.Price).Quo(sdk.NewDec(int64(math.Pow10(int(msg.PricePrecision))))),
		Quantity:    msg.Quantity,
		Side:        msg.Side,
		TimeInForce: msg.TimeInForce,
		Height:      ctx.BlockHeight(),
		ExistBlocks: existBlocks,
		FrozenFee:   frozenFee,
		LeftStock:   msg.Quantity,
		Freeze:      amount,
		DealMoney:   0,
		DealStock:   0,
	}

	ork := keepers.NewOrderKeeper(keeper.GetMarketKey(), order.TradingPair, types.ModuleCdc)
	if err := ork.Add(ctx, &order); err != nil {
		return err.Result()
	}

	if err := handleFeeForCreateOrder(ctx, keeper, amount, denom, order.Sender, frozenFee, featureFee); err != nil {
		return err.Result()
	}
	sendCreateOrderMsg(ctx, keeper, order, featureFee)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyOrder, order.OrderID())),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyTradingPair, order.TradingPair),
			sdk.NewAttribute(AttributeKeyHeight, strconv.FormatInt(order.Height, 10))),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCreateOrder(ctx sdk.Context, keeper keepers.Keeper, msg types.MsgCreateOrder, cetFee int64, amount int64, denom string, seq uint64) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	if cetFee != 0 {
		if !keeper.HasCoins(ctx, msg.Sender, sdk.Coins{sdk.NewCoin(dex.CET, sdk.NewInt(cetFee))}) {
			return types.ErrInsufficientCoins().Result()
		}
	}
	stock, money := SplitSymbol(msg.TradingPair)
	totalAmount := sdk.NewInt(amount)
	if (stock == dex.CET && msg.Side == types.SELL) ||
		(money == dex.CET && msg.Side == types.BUY) {
		totalAmount = totalAmount.AddRaw(cetFee)
	}
	if !keeper.HasCoins(ctx, msg.Sender, sdk.Coins{sdk.NewCoin(denom, totalAmount)}) {
		return types.ErrInsufficientCoins().Result()
	}
	orderID, err := types.AssemblyOrderID(msg.Sender.String(), seq, msg.Identify)
	if err != nil {
		return types.ErrInvalidSequence(err.Error()).Result()
	}
	globalKeeper := keepers.NewGlobalOrderKeeper(keeper.GetMarketKey(), types.ModuleCdc)
	if globalKeeper.QueryOrder(ctx, orderID) != nil {
		return types.ErrOrderAlreadyExist(orderID).Result()
	}
	marketInfo, err := keeper.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return types.ErrInvalidMarket(err.Error()).Result()
	}
	if p := msg.PricePrecision; p > marketInfo.PricePrecision {
		return types.ErrInvalidPricePrecision(p).Result()
	}
	if keeper.IsTokenForbidden(ctx, stock) || keeper.IsTokenForbidden(ctx, money) {
		return types.ErrTokenForbidByIssuer().Result()
	}
	if keeper.IsForbiddenByTokenIssuer(ctx, stock, msg.Sender) || keeper.IsForbiddenByTokenIssuer(ctx, money, msg.Sender) {
		return types.ErrAddressForbidByIssuer().Result()
	}
	baseValue := types.GetGranularityOfOrder(marketInfo.OrderPrecision)
	if amount%baseValue != 0 {
		return types.ErrInvalidOrderAmount("The amount of tokens to trade should be a multiple of the order precision").Result()
	}

	return sdk.Result{}
}

func handleMsgCancelOrder(ctx sdk.Context, msg types.MsgCancelOrder, keeper keepers.Keeper) sdk.Result {
	if err := checkMsgCancelOrder(ctx, msg, keeper); err != nil {
		return err.Result()
	}

	order := keepers.NewGlobalOrderKeeper(keeper.GetMarketKey(), types.ModuleCdc).QueryOrder(ctx, msg.OrderID)
	marketParams := keeper.GetParams(ctx)

	ork := keepers.NewOrderKeeper(keeper.GetMarketKey(), order.TradingPair, types.ModuleCdc)
	removeOrder(ctx, ork, keeper.GetBankxKeeper(), keeper, order, marketParams.FeeForZeroDeal)

	// send msg to kafka
	msgInfo := types.CancelOrderInfo{
		OrderID:        msg.OrderID,
		TradingPair:    order.TradingPair,
		Height:         ctx.BlockHeight(),
		Side:           order.Side,
		Price:          order.Price,
		DelReason:      types.CancelOrderByManual,
		UsedCommission: order.CalOrderFeeInt64(marketParams.FeeForZeroDeal),
		LeftStock:      order.LeftStock,
		RemainAmount:   order.Freeze,
		DealStock:      order.DealStock,
		DealMoney:      order.DealMoney,
	}
	fillMsgQueue(ctx, keeper, types.CancelOrderInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyOrder, order.OrderID())),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeKeyDelOrderReason, types.CancelOrderByManual),
			sdk.NewAttribute(AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyDelOrderHeight, strconv.Itoa(int(ctx.BlockHeight()))),
			sdk.NewAttribute(AttributeKeyTradingPair, order.TradingPair),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCancelOrder(ctx sdk.Context, msg types.MsgCancelOrder, keeper keepers.Keeper) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	globalKeeper := keepers.NewGlobalOrderKeeper(keeper.GetMarketKey(), types.ModuleCdc)
	order := globalKeeper.QueryOrder(ctx, msg.OrderID)
	if order == nil {
		return types.ErrOrderNotFound(msg.OrderID)
	}

	if !bytes.Equal(order.Sender, msg.Sender) {
		return types.ErrNotMatchSender("only order's sender can cancel this order")
	}

	return nil
}

func handleMsgCancelTradingPair(ctx sdk.Context, msg types.MsgCancelTradingPair, keeper keepers.Keeper) sdk.Result {
	if err := checkMsgCancelTradingPair(keeper, msg, ctx); err != nil {
		return err.Result()
	}

	// Add del request to store
	dlk := keepers.NewDelistKeeper(keeper.GetMarketKey())
	delistSymbols := dlk.GetDelistSymbolsBeforeTime(ctx, math.MaxInt64)
	for _, sym := range delistSymbols {
		if msg.TradingPair == sym {
			return types.ErrDelistRequestExist(sym).Result()
		}
	}
	dlk.AddDelistRequest(ctx, msg.EffectiveTime, msg.TradingPair)

	// send msg to kafka
	//values := strings.Split(msg.TradingPair, types.SymbolSeparator)
	//msgInfo := types.CancelMarketInfo{
	//	Stock:   values[0],
	//	Money:   values[1],
	//	Deleter: msg.Sender.String(),
	//	DelTime: msg.EffectiveTime,
	//}
	//fillMsgQueue(ctx, keeper, types.CancelTradingInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyTradingPair, msg.TradingPair)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeKeyEffectiveTime, strconv.Itoa(int(msg.EffectiveTime))),
			sdk.NewAttribute(AttributeKeySender, msg.Sender.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCancelTradingPair(keeper keepers.Keeper, msg types.MsgCancelTradingPair, ctx sdk.Context) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	marketParams := keeper.GetParams(ctx)
	currTime := ctx.BlockHeader().Time.Unix()
	if msg.EffectiveTime < currTime+marketParams.MarketMinExpiredTime {
		return types.ErrInvalidCancelTime()
	}

	info, err := keeper.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return types.ErrInvalidMarket(err.Error())
	}

	stockToken := keeper.GetToken(ctx, info.Stock)
	if !bytes.Equal(msg.Sender, stockToken.GetOwner()) {
		return types.ErrNotMatchSender("only stock's owner can cancel a market")
	}

	// TODO. Will add unit test
	if !stockToken.GetTokenForbiddable() {
		if info.Money == dex.CET {
			return types.ErrDelistNotAllowed("stock token doesn't have globally forbidden attribute, so its market against CET can not be canceled")
		}
	}

	// TODO. Will add unit test
	if info.Money == dex.CET && keeper.IsBancorExist(ctx, info.Stock) {
		return types.ErrDelistNotAllowed(
			fmt.Sprintf("When %s has bancor contracts, you can't delist the %s/cet market", info.Stock, info.Stock))
	}

	return nil
}

func calculateAmount(price, quantity int64, pricePrecision byte) (sdk.Dec, error) {
	actualPrice := sdk.NewDec(price).Quo(sdk.NewDec(int64(math.Pow10(int(pricePrecision)))))
	money := actualPrice.Mul(sdk.NewDec(quantity)).Add(sdk.NewDec(types.ExtraFrozenMoney)).Ceil()
	if money.GT(sdk.NewDec(types.MaxOrderAmount)) {
		return money, fmt.Errorf("exchange amount exceeds max int64 ")
	}
	return money, nil
}

func handleMsgModifyPricePrecision(ctx sdk.Context, msg types.MsgModifyPricePrecision, k keepers.Keeper) sdk.Result {
	if err := checkMsgModifyPricePrecision(ctx, msg, k); err != nil {
		return err.Result()
	}

	oldInfo, _ := k.GetMarketInfo(ctx, msg.TradingPair)
	info := types.MarketInfo{
		Stock:             oldInfo.Stock,
		Money:             oldInfo.Money,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: oldInfo.LastExecutedPrice,
	}
	if err := k.SetMarket(ctx, info); err != nil {
		return err.Result()
	}

	//msgInfo := types.ModifyPricePrecisionInfo{
	//	Sender:            msg.Sender.String(),
	//	TradingPair:       msg.TradingPair,
	//	OldPricePrecision: oldInfo.PricePrecision,
	//	NewPricePrecision: info.PricePrecision,
	//}
	//fillMsgQueue(ctx, k, types.PricePrecisionInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyTradingPair, msg.TradingPair)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyOldPricePrecision, strconv.Itoa(int(oldInfo.PricePrecision))),
			sdk.NewAttribute(AttributeKeyNewPricePrecision, strconv.Itoa(int(info.PricePrecision)))),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgModifyPricePrecision(ctx sdk.Context, msg types.MsgModifyPricePrecision, k keepers.Keeper) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	info, err := k.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return types.ErrInvalidMarket("Error retrieving market information: " + err.Error())
	}

	if info.PricePrecision > msg.PricePrecision {
		return types.ErrInvalidPricePrecisionChange(fmt.Sprintf(
			"Price Precision can only be increased; tradingPair price_precision : %d, msg price_precision : %d",
			info.PricePrecision, msg.PricePrecision))
	}

	stock, _ := SplitSymbol(msg.TradingPair)
	tokenInfo := k.GetToken(ctx, stock)
	if !tokenInfo.GetOwner().Equals(msg.Sender) {
		return types.ErrNotMatchSender(fmt.Sprintf(
			"The sender of the transaction (%s) does not match the owner of the transaction pair (%s)",
			tokenInfo.GetOwner().String(), msg.Sender.String()))
	}

	return nil
}
