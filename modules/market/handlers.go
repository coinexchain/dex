package market

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	mtype "github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k keepers.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case mtype.MsgCreateTradingPair:
			return handleMsgCreateTradingPair(ctx, msg, k)
		case mtype.MsgCreateOrder:
			return handleMsgCreateOrder(ctx, msg, k)
		case mtype.MsgCancelOrder:
			return handleMsgCancelOrder(ctx, msg, k)
		case mtype.MsgCancelTradingPair:
			return handleMsgCancelTradingPair(ctx, msg, k)
		case mtype.MsgModifyPricePrecision:
			return handleMsgModifyPricePrecision(ctx, msg, k)
		default:
			errMsg := "Unrecognized market Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateTradingPair(ctx sdk.Context, msg mtype.MsgCreateTradingPair, keeper keepers.Keeper) sdk.Result {
	if ret := checkMsgCreateTradingPair(ctx, msg, keeper); !ret.IsOK() {
		return ret
	}

	info := mtype.MarketInfo{
		Stock:             msg.Stock,
		Money:             msg.Money,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: sdk.ZeroDec(),
	}

	if err := keeper.SetMarket(ctx, info); err != nil {
		return err.Result()
	}

	param := keeper.GetParams(ctx)
	if err := keeper.SubtractFeeAndCollectFee(ctx, msg.Creator, types.NewCetCoins(param.CreateMarketFee)); err != nil {
		return err.Result()
	}

	// send msg to kafka
	msgInfo := mtype.CreateMarketInfo{
		Stock:          msg.Stock,
		Money:          msg.Money,
		PricePrecision: msg.PricePrecision,
		Creator:        msg.Creator.String(),
		CreateHeight:   ctx.BlockHeight(),
	}
	fillMsgQueue(ctx, keeper, mtype.CreateTradingInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeMarket,
			sdk.NewAttribute(AttributeKeyTradingPair, msg.Stock+mtype.SymbolSeparator+msg.Money),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, mtype.ModuleName),
			sdk.NewAttribute(AttributeKeyStock, msg.Stock),
			sdk.NewAttribute(AttributeKeyMoney, msg.Money),
			sdk.NewAttribute(AttributeKeyPricePrecision, strconv.Itoa(int(info.PricePrecision))),
			sdk.NewAttribute(AttributeKeyLastExecutePrice, info.LastExecutedPrice.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper keepers.Keeper, key string, msg interface{}) {
	if keeper.IsSubScribe(mtype.Topic) {
		fillMsgs(ctx, key, msg)
	}
}

func fillMsgs(ctx sdk.Context, key string, msg interface{}) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(msgqueue.EventTypeMsgQueue,
		sdk.NewAttribute(key, string(bytes))))
}

func checkMsgCreateTradingPair(ctx sdk.Context, msg mtype.MsgCreateTradingPair, keeper keepers.Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if _, err := keeper.GetMarketInfo(ctx, msg.Stock+mtype.SymbolSeparator+msg.Money); err == nil {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeRepeatTrade, "The repeatedly created trading pairs").Result()
	}

	if !keeper.IsTokenExists(ctx, msg.Money) || !keeper.IsTokenExists(ctx, msg.Stock) {
		return mtype.ErrTokenNoExist().Result()
	}

	if !keeper.IsTokenIssuer(ctx, msg.Stock, []byte(msg.Creator)) {
		return mtype.ErrInvalidTokenIssuer().Result()
	}

	if msg.Money != types.CET && msg.Stock != types.CET {
		if _, err := keeper.GetMarketInfo(ctx, msg.Stock+mtype.SymbolSeparator+types.CET); err != nil {
			return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeStockNoHaveCetTrade, "The stock(%s) not have cet trade", msg.Stock).Result()
		}
	}

	if msg.PricePrecision < mtype.MinTokenPricePrecision || msg.PricePrecision > mtype.MaxTokenPricePrecision {
		return mtype.ErrInvalidPricePrecision().Result()
	}

	marketParams := keeper.GetParams(ctx)
	if !keeper.HasCoins(ctx, msg.Creator, types.NewCetCoins(marketParams.CreateMarketFee)) {
		return mtype.ErrInsufficientCoins().Result()
	}

	return sdk.Result{}
}

func calFrozenFeeInOrder(ctx sdk.Context, marketParams keepers.Params, keeper keepers.Keeper, msg mtype.MsgCreateOrder) (int64, sdk.Error) {
	var frozenFee int64
	stock := strings.Split(msg.TradingPair, mtype.SymbolSeparator)[0]

	// Calculate the fee when stock is cet
	rate := sdk.NewDec(marketParams.MarketFeeRate)
	div := sdk.NewDec(int64(math.Pow10(keepers.MarketFeeRatePrecision)))
	if stock == types.CET {
		frozenFee = sdk.NewDec(msg.Quantity).Mul(rate).Quo(div).RoundInt64()
	} else {
		stockSepCet := stock + mtype.SymbolSeparator + types.CET
		marketInfo, err := keeper.GetMarketInfo(ctx, stockSepCet)
		if err != nil || marketInfo.LastExecutedPrice.IsZero() {
			frozenFee = marketParams.FixedTradeFee
		} else {
			totalPriceInCet := marketInfo.LastExecutedPrice.Mul(sdk.NewDec(msg.Quantity))
			frozenFee = totalPriceInCet.Mul(rate).Quo(div).RoundInt64()
		}
	}
	if frozenFee < marketParams.MarketFeeMin {
		return 0, mtype.ErrOrderQuantityToSmall()
	}

	return frozenFee, nil
}

func calFeatureFeeForExistBlocks(msg mtype.MsgCreateOrder, marketParam keepers.Params) int64 {
	if msg.TimeInForce == mtype.IOC {
		return 0
	}
	quotient := msg.ExistBlocks / marketParam.GTEOrderLifetime
	remainder := msg.ExistBlocks % marketParam.GTEOrderLifetime
	if remainder != 0 {
		quotient++
	}
	if quotient <= 1 {
		return 0
	}
	return int64(quotient) * marketParam.GTEOrderFeatureFeeByBlocks
}

func handleFeeForCreateOrder(ctx sdk.Context, keeper keepers.Keeper, amount int64, denom string,
	sender sdk.AccAddress, frozenFee, featureFee int64) sdk.Error {
	coin := sdk.NewCoin(denom, sdk.NewInt(amount))
	if err := keeper.FreezeCoins(ctx, sender, sdk.Coins{coin}); err != nil {
		return err
	}
	if frozenFee != 0 {
		frozenFeeAsCet := sdk.Coins{sdk.NewCoin(types.CET, sdk.NewInt(frozenFee))}
		if err := keeper.FreezeCoins(ctx, sender, frozenFeeAsCet); err != nil {
			return err
		}
	}
	if featureFee != 0 {
		if err := keeper.SubtractFeeAndCollectFee(ctx, sender, types.NewCetCoins(featureFee)); err != nil {
			return err
		}
	}
	return nil
}

func sendCreateOrderMsg(ctx sdk.Context, keeper keepers.Keeper, order mtype.Order, featureFee int64) {
	// send msg to kafka
	msgInfo := mtype.CreateOrderInfo{
		OrderID:     order.OrderID(),
		Sender:      order.Sender.String(),
		TradingPair: order.TradingPair,
		OrderType:   order.OrderType,
		Price:       order.Price.String(),
		Quantity:    order.Quantity,
		Side:        order.Side,
		TimeInForce: order.TimeInForce,
		Height:      order.Height,
		FrozenFee:   order.FrozenFee,
		Freeze:      order.Freeze,
		FeatureFee:  featureFee,
	}
	fillMsgQueue(ctx, keeper, mtype.CreateOrderInfoKey, msgInfo)
}

func handleMsgCreateOrder(ctx sdk.Context, msg mtype.MsgCreateOrder, keeper keepers.Keeper) sdk.Result {

	values := strings.Split(msg.TradingPair, mtype.SymbolSeparator)
	stock, money := values[0], values[1]
	denom := stock
	amount := msg.Quantity
	if msg.Side == mtype.BUY {
		denom = money
		amount = calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision).RoundInt64()
	}
	if amount > mtype.MaxOrderAmount {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidOrderAmount, "The order amount is too large").Result()
	}

	marketParams := keeper.GetParams(ctx)
	frozenFee, err := calFrozenFeeInOrder(ctx, marketParams, keeper, msg)
	if err != nil {
		return err.Result()
	}

	featureFee := calFeatureFeeForExistBlocks(msg, marketParams)
	if ret := checkMsgCreateOrder(ctx, keeper, msg, frozenFee+featureFee, amount, denom); !ret.IsOK() {
		return ret
	}

	order := mtype.Order{
		Sender:      msg.Sender,
		Sequence:    msg.Sequence,
		TradingPair: msg.TradingPair,
		OrderType:   msg.OrderType,
		Price:       sdk.NewDec(msg.Price).Quo(sdk.NewDec(int64(math.Pow10(int(msg.PricePrecision))))),
		Quantity:    msg.Quantity,
		Side:        msg.Side,
		TimeInForce: msg.TimeInForce,
		Height:      ctx.BlockHeight(),
		ExistBlocks: msg.ExistBlocks,
		FrozenFee:   frozenFee,
		LeftStock:   msg.Quantity,
		Freeze:      amount,
		DealMoney:   0,
		DealStock:   0,
	}

	ork := keepers.NewOrderKeeper(keeper.GetMarketKey(), order.TradingPair, mtype.ModuleCdc)
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
			sdk.NewAttribute(sdk.AttributeKeyModule, mtype.ModuleName),
			sdk.NewAttribute(AttributeKeyTradingPair, order.TradingPair),
			sdk.NewAttribute(AttributeKeyHeight, strconv.FormatInt(order.Height, 10)),
			sdk.NewAttribute(AttributeKeySequence, strconv.FormatInt(int64(order.Sequence), 10)),
			sdk.NewAttribute(AttributeKeyOrderType, strconv.Itoa(int(order.OrderType))),
			sdk.NewAttribute(AttributeKeySide, strconv.Itoa(int(order.Side))),
			sdk.NewAttribute(AttributeKeyPrice, order.Price.String()),
			sdk.NewAttribute(AttributeKeyQuantity, strconv.FormatInt(order.Quantity, 10)),
			sdk.NewAttribute(AttributeKeyTimeInForce, strconv.Itoa(order.TimeInForce))),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCreateOrder(ctx sdk.Context, keeper keepers.Keeper, msg mtype.MsgCreateOrder, cetFee int64, amount int64, denom string) sdk.Result {
	var (
		stock, money string
	)
	values := strings.Split(msg.TradingPair, mtype.SymbolSeparator)
	stock, money = values[0], values[1]

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if cetFee != 0 {
		frozenFeeAsCet := sdk.Coins{sdk.NewCoin(types.CET, sdk.NewInt(cetFee))}
		if !keeper.HasCoins(ctx, msg.Sender, frozenFeeAsCet) {
			return mtype.ErrInsufficientCoins().Result()
		}
	}

	totalAmount := amount
	if (stock == types.CET && msg.Side == mtype.SELL) ||
		(money == types.CET && msg.Side == mtype.BUY) {
		totalAmount += cetFee
	}
	if !keeper.HasCoins(ctx, msg.Sender, sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(totalAmount))}) {
		return mtype.ErrInsufficientCoins().Result()
	}

	marketInfo, err := keeper.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return mtype.ErrInvalidSymbol().Result()
	}
	if msg.PricePrecision > marketInfo.PricePrecision {
		return mtype.ErrInvalidPricePrecision().Result()
	}

	if keeper.IsTokenForbidden(ctx, stock) || keeper.IsTokenForbidden(ctx, money) {
		return mtype.ErrTokenForbidByIssuer().Result()
	}

	if keeper.IsForbiddenByTokenIssuer(ctx, stock, msg.Sender) || keeper.IsForbiddenByTokenIssuer(ctx, money, msg.Sender) {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeAddressForbidByIssuer, "The sender is forbidden by token issuer").Result()
	}

	return sdk.Result{}
}

func handleMsgCancelOrder(ctx sdk.Context, msg mtype.MsgCancelOrder, keeper keepers.Keeper) sdk.Result {

	if err := checkMsgCancelOrder(ctx, msg, keeper); !err.IsOK() {
		return err
	}

	order := keepers.NewGlobalOrderKeeper(keeper.GetMarketKey(), mtype.ModuleCdc).QueryOrder(ctx, msg.OrderID)
	marketParams := keeper.GetParams(ctx)

	ork := keepers.NewOrderKeeper(keeper.GetMarketKey(), order.TradingPair, mtype.ModuleCdc)
	removeOrder(ctx, ork, keeper.GetBankxKeeper(), keeper, order, marketParams.FeeForZeroDeal)

	// send msg to kafka
	msgInfo := mtype.CancelOrderInfo{
		OrderID:        msg.OrderID,
		DelReason:      mtype.CancelOrderByManual,
		DelHeight:      ctx.BlockHeight(),
		UsedCommission: order.CalOrderFee(marketParams.FeeForZeroDeal).RoundInt64(),
		LeftStock:      order.LeftStock,
		RemainAmount:   order.Freeze,
		DealStock:      order.DealStock,
		DealMoney:      order.DealMoney,
	}
	fillMsgQueue(ctx, keeper, mtype.CancelOrderInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyOrder, order.OrderID())),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, mtype.ModuleName),
			sdk.NewAttribute(AttributeKeyDelOrderReason, mtype.CancelOrderByManual),
			sdk.NewAttribute(AttributeKeyDelOrderHeight, strconv.Itoa(int(ctx.BlockHeight()))),
			sdk.NewAttribute(AttributeKeyUsedCommission, order.CalOrderFee(marketParams.FeeForZeroDeal).String()),
			sdk.NewAttribute(AttributeKeyLeftStock, strconv.Itoa(int(order.LeftStock))),
			sdk.NewAttribute(AttributeKeyDealStock, strconv.Itoa(int(order.DealStock))),
			sdk.NewAttribute(AttributeKeyRemainAmount, strconv.Itoa(int(order.FrozenFee))),
			sdk.NewAttribute(AttributeKeyDealMoney, strconv.Itoa(int(order.DealMoney))),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCancelOrder(ctx sdk.Context, msg mtype.MsgCancelOrder, keeper keepers.Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	globalKeeper := keepers.NewGlobalOrderKeeper(keeper.GetMarketKey(), mtype.ModuleCdc)
	order := globalKeeper.QueryOrder(ctx, msg.OrderID)
	if order == nil {
		return sdk.NewError(mtype.StoreKey, mtype.CodeNotFindOrder, "Not find order in blockchain").Result()
	}

	if !bytes.Equal(order.Sender, msg.Sender) {
		return sdk.NewError(mtype.StoreKey, mtype.CodeNotMatchSender, "The cancel addr is not match order sender").Result()
	}

	return sdk.Result{}
}

func handleMsgCancelTradingPair(ctx sdk.Context, msg mtype.MsgCancelTradingPair, keeper keepers.Keeper) sdk.Result {

	if err := checkMsgCancelTradingPair(keeper, msg, ctx); err != nil {
		return err.Result()
	}

	// Add del request to store
	dlk := keepers.NewDelistKeeper(keeper.GetMarketKey())
	dlk.AddDelistRequest(ctx, msg.EffectiveTime, msg.TradingPair)

	// send msg to kafka
	values := strings.Split(msg.TradingPair, mtype.SymbolSeparator)
	msgInfo := mtype.CancelMarketInfo{
		Stock:   values[0],
		Money:   values[1],
		Deleter: msg.Sender.String(),
		DelTime: msg.EffectiveTime,
	}
	fillMsgQueue(ctx, keeper, mtype.CancelTradingInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyTradingPair, msg.TradingPair)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, mtype.ModuleName),
			sdk.NewAttribute(AttributeKeyEffectiveTime, strconv.Itoa(int(msg.EffectiveTime)))),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgCancelTradingPair(keeper keepers.Keeper, msg mtype.MsgCancelTradingPair, ctx sdk.Context) sdk.Error {

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	marketParams := keeper.GetParams(ctx)
	currTime := ctx.BlockHeader().Time.Unix()
	if msg.EffectiveTime < currTime+marketParams.MarketMinExpiredTime {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidTime, "Invalid Cancel Time")
	}

	info, err := keeper.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidSymbol, err.Error())
	}

	stockToken := keeper.GetToken(ctx, info.Stock)
	if !bytes.Equal(msg.Sender, stockToken.GetOwner()) {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeNotMatchSender, "Not match market info sender")
	}

	// TODO. Will add unit test
	if !stockToken.GetTokenForbiddable() {
		if info.Money == types.CET {
			return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeNotAllowedOffline, "stock token don't have globally forbidden attribute, so a trade with CET would not be allowed ")
		}
	}

	if keeper.IsBancorExist(ctx, info.Stock) {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidBancorExist,
			"When stock has bancor contracts, you can't delete the trading-pair")
	}

	return nil
}

func calculateAmount(price, quantity int64, pricePrecision byte) sdk.Dec {
	actualPrice := sdk.NewDec(price).Quo(sdk.NewDec(int64(math.Pow10(int(pricePrecision)))))
	money := actualPrice.Mul(sdk.NewDec(quantity))
	return money.Add(sdk.NewDec(mtype.ExtraFrozenMoney)).Ceil()
}

func handleMsgModifyPricePrecision(ctx sdk.Context, msg mtype.MsgModifyPricePrecision, k keepers.Keeper) sdk.Result {
	if ret := checkMsgModifyPricePrecision(ctx, msg, k); !ret.IsOK() {
		return ret
	}

	oldInfo, _ := k.GetMarketInfo(ctx, msg.TradingPair)
	info := mtype.MarketInfo{
		Stock:             oldInfo.Stock,
		Money:             oldInfo.Money,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: oldInfo.LastExecutedPrice,
	}
	if err := k.SetMarket(ctx, info); err != nil {
		return err.Result()
	}

	msgInfo := mtype.ModifyPricePrecisionInfo{
		Sender:            msg.Sender.String(),
		TradingPair:       msg.TradingPair,
		OldPricePrecision: oldInfo.PricePrecision,
		NewPricePrecision: info.PricePrecision,
	}
	fillMsgQueue(ctx, k, mtype.PricePrecisionInfoKey, msgInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(EventTypeMarket, sdk.NewAttribute(
			AttributeKeyTradingPair, msg.TradingPair)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, mtype.ModuleName),
			sdk.NewAttribute(AttributeKeyOldPricePrecision, strconv.Itoa(int(oldInfo.PricePrecision))),
			sdk.NewAttribute(AttributeKeyNewPricePrecision, strconv.Itoa(int(info.PricePrecision)))),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func checkMsgModifyPricePrecision(ctx sdk.Context, msg mtype.MsgModifyPricePrecision, k keepers.Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	info, err := k.GetMarketInfo(ctx, msg.TradingPair)
	if err != nil {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidSymbol,
			fmt.Sprintf("Error retrieving trade pair information : %s", err.Error())).Result()
	}

	if info.PricePrecision > msg.PricePrecision {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeInvalidPricePrecision,
			fmt.Sprintf("Price Precision can only be increased; "+
				"tradingPair price_precision : %d, msg price_precision : %d",
				&info.PricePrecision, msg.PricePrecision)).Result()
	}

	stock := strings.Split(msg.TradingPair, mtype.SymbolSeparator)[0]
	tokenInfo := k.GetToken(ctx, stock)
	if !tokenInfo.GetOwner().Equals(msg.Sender) {
		return sdk.NewError(mtype.CodeSpaceMarket, mtype.CodeNotMatchSender, fmt.Sprintf(
			"The sender of the transaction (%s) does not match "+
				"the owner of the transaction pair (%s)",
			tokenInfo.GetOwner().String(), msg.Sender.String())).Result()
	}

	return sdk.Result{}
}
