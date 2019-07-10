package market

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/coinexchain/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/match"
)

const (
	MinTokenPricePrecision           = 8
	MaxTokenPricePrecision           = 18
	LimitOrder             OrderType = 2
	SymbolSeparator                  = "/"
	ExtraFrozenMoney                 = 0 // 100
)

type OrderType = byte

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateTradingPair:
			return handleMsgCreateTradingPair(ctx, msg, k)
		case MsgCreateOrder:
			return handleMsgCreateOrder(ctx, msg, k)
		case MsgCancelOrder:
			return handleMsgCancelOrder(ctx, msg, k)
		case MsgCancelTradingPair:
			return handleMsgCancelTradingPair(ctx, msg, k)
		default:
			errMsg := "Unrecognized market Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateTradingPair(ctx sdk.Context, msg MsgCreateTradingPair, keeper Keeper) sdk.Result {
	if ret := checkMsgCreateTradingPair(ctx, msg, keeper); !ret.IsOK() {
		return ret
	}

	info := MarketInfo{
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
		// Here must panic. because the market info have stored in db.
		panic(err)
	}

	// send msg to kafka
	msgInfo := CreateMarketInfo{
		Stock:          msg.Stock,
		Money:          msg.Money,
		PricePrecision: msg.PricePrecision,
		Creator:        msg.Creator.String(),
		CreateHeight:   ctx.BlockHeight(),
	}
	keeper.msgProducer.SendMsg(Topic, CreateMarketInfoKey, msgInfo)
	return sdk.Result{Tags: info.GetTags()}
}

func checkMsgCreateTradingPair(ctx sdk.Context, msg MsgCreateTradingPair, keeper Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if _, err := keeper.GetMarketInfo(ctx, msg.Stock+SymbolSeparator+msg.Money); err == nil {
		return sdk.NewError(CodeSpaceMarket, CodeRepeatTrade, "The repeatedly created trading pairs").Result()
	}

	if !keeper.axk.IsTokenExists(ctx, msg.Money) || !keeper.axk.IsTokenExists(ctx, msg.Stock) {
		return ErrTokenNoExist().Result()
	}

	if !keeper.axk.IsTokenIssuer(ctx, msg.Stock, []byte(msg.Creator)) {
		return ErrInvalidTokenIssuer().Result()
	}

	if msg.Money != types.CET && msg.Stock != types.CET {
		if _, err := keeper.GetMarketInfo(ctx, msg.Stock+SymbolSeparator+types.CET); err != nil {
			return sdk.NewError(CodeSpaceMarket, CodeStockNoHaveCetTrade, "The stock(%s) not have cet trade", msg.Stock).Result()
		}
	}

	if msg.PricePrecision < MinTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	marketParams := keeper.GetParams(ctx)
	if !keeper.bnk.HasCoins(ctx, msg.Creator, types.NewCetCoins(marketParams.CreateMarketFee)) {
		return ErrInsufficientCoins().Result()
	}

	return sdk.Result{}
}

func calFrozenFeeInOrder(ctx sdk.Context, marketParams Params, keeper Keeper, msg MsgCreateOrder) (int64, sdk.Error) {
	var frozenFee int64
	stock := strings.Split(msg.Symbol, SymbolSeparator)[0]

	// Calculate the fee when stock is cet
	rate := sdk.NewDec(marketParams.MarketFeeRate)
	div := sdk.NewDec(int64(math.Pow10(MarketFeeRatePrecision)))
	if stock == types.CET {
		frozenFee = sdk.NewDec(msg.Quantity).Mul(rate).Quo(div).RoundInt64()
	} else {
		stockSepCet := stock + SymbolSeparator + types.CET
		marketInfo, err := keeper.GetMarketInfo(ctx, stockSepCet)
		if err != nil || marketInfo.LastExecutedPrice.IsZero() {
			frozenFee = marketParams.FixedTradeFee
		} else {
			totalPriceInCet := marketInfo.LastExecutedPrice.Mul(sdk.NewDec(msg.Quantity))
			frozenFee = totalPriceInCet.Mul(rate).Quo(div).RoundInt64()
		}
	}
	if frozenFee < marketParams.MarketFeeMin {
		return 0, ErrOrderQuantityToSmall()
	}

	return frozenFee, nil
}

func calFeatureFeeForExistBlocks(msg MsgCreateOrder, marketParam Params) int64 {
	if msg.TimeInForce == IOC {
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

func handleFeeForCreateOrder(ctx sdk.Context, keeper Keeper, amount int64, denom string,
	sender sdk.AccAddress, frozenFee, featureFee int64) sdk.Error {
	coin := sdk.NewCoin(denom, sdk.NewInt(amount))
	if err := keeper.bnk.FreezeCoins(ctx, sender, sdk.Coins{coin}); err != nil {
		return err
	}
	if frozenFee != 0 {
		frozenFeeAsCet := sdk.Coins{sdk.NewCoin(types.CET, sdk.NewInt(frozenFee))}
		if err := keeper.bnk.FreezeCoins(ctx, sender, frozenFeeAsCet); err != nil {
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

func sendCreateOrderMsg(keeper Keeper, order Order, featureFee int64) {
	// send msg to kafka
	msgInfo := CreateOrderInfo{
		OrderID:     order.OrderID(),
		Sender:      order.Sender.String(),
		Symbol:      order.Symbol,
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
	keeper.SendMsg(CreateOrderInfoKey, msgInfo)
}

func handleMsgCreateOrder(ctx sdk.Context, msg MsgCreateOrder, keeper Keeper) sdk.Result {

	values := strings.Split(msg.Symbol, SymbolSeparator)
	stock, money := values[0], values[1]
	denom := stock
	amount := msg.Quantity
	if msg.Side == match.BUY {
		denom = money
		amount = calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision).RoundInt64()
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

	order := Order{
		Sender:      msg.Sender,
		Sequence:    msg.Sequence,
		Symbol:      msg.Symbol,
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

	ork := NewOrderKeeper(keeper.marketKey, order.Symbol, keeper.cdc)
	if err := ork.Add(ctx, &order); err != nil {
		return err.Result()
	}

	if err := handleFeeForCreateOrder(ctx, keeper, amount, denom, order.Sender, frozenFee, featureFee); err != nil {
		return err.Result()
	}
	sendCreateOrderMsg(keeper, order, featureFee)

	return sdk.Result{Tags: order.GetTagsInOrderCreate()}
}

func checkMsgCreateOrder(ctx sdk.Context, keeper Keeper, msg MsgCreateOrder, cetFee int64, amount int64, denom string) sdk.Result {
	var (
		stock, money string
	)
	values := strings.Split(msg.Symbol, SymbolSeparator)
	stock, money = values[0], values[1]

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if cetFee != 0 {
		frozenFeeAsCet := sdk.Coins{sdk.NewCoin(types.CET, sdk.NewInt(cetFee))}
		if !keeper.bnk.HasCoins(ctx, msg.Sender, frozenFeeAsCet) {
			return ErrInsufficientCoins().Result()
		}
	}

	totalAmount := amount
	if (stock == types.CET && msg.Side == match.SELL) ||
		(money == types.CET && msg.Side == match.BUY) {
		totalAmount += cetFee
	}
	if !keeper.bnk.HasCoins(ctx, msg.Sender, sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(totalAmount))}) {
		return ErrInsufficientCoins().Result()
	}

	marketInfo, err := keeper.GetMarketInfo(ctx, msg.Symbol)
	if err != nil {
		return ErrInvalidSymbol().Result()
	}
	if msg.PricePrecision > marketInfo.PricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	if keeper.axk.IsTokenForbidden(ctx, stock) || keeper.axk.IsTokenForbidden(ctx, money) {
		return ErrTokenForbidByIssuer().Result()
	}

	if keeper.axk.IsForbiddenByTokenIssuer(ctx, stock, msg.Sender) || keeper.axk.IsForbiddenByTokenIssuer(ctx, money, msg.Sender) {
		return sdk.NewError(CodeSpaceMarket, CodeAddressForbidByIssuer, "The sender is forbidden by token issuer").Result()
	}

	return sdk.Result{}
}

func handleMsgCancelOrder(ctx sdk.Context, msg MsgCancelOrder, keeper Keeper) sdk.Result {

	if err := checkMsgCancelOrder(ctx, msg, keeper); !err.IsOK() {
		return err
	}

	order := NewGlobalOrderKeeper(keeper.marketKey, keeper.cdc).QueryOrder(ctx, msg.OrderID)
	marketParams := keeper.GetParams(ctx)

	ork := NewOrderKeeper(keeper.marketKey, order.Symbol, keeper.cdc)
	removeOrder(ctx, ork, keeper.bnk, keeper, order, marketParams.FeeForZeroDeal)

	// send msg to kafka
	msgInfo := CancelOrderInfo{
		OrderID:        msg.OrderID,
		DelReason:      CancelOrderByManual,
		DelHeight:      ctx.BlockHeight(),
		UsedCommission: order.CalOrderFee(marketParams.FeeForZeroDeal).RoundInt64(),
		LeftStock:      order.LeftStock,
		RemainAmount:   order.Freeze,
		DealStock:      order.DealStock,
		DealMoney:      order.DealMoney,
	}
	keeper.SendMsg(CancelOrderInfoKey, msgInfo)

	return sdk.Result{Tags: sdk.NewTags(
		Sender, msg.Sender.String(),
		OrderID, msg.OrderID,
	),
	}
}

func checkMsgCancelOrder(ctx sdk.Context, msg MsgCancelOrder, keeper Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	globalKeeper := NewGlobalOrderKeeper(keeper.marketKey, keeper.cdc)
	order := globalKeeper.QueryOrder(ctx, msg.OrderID)
	if order == nil {
		return sdk.NewError(StoreKey, CodeNotFindOrder, "Not find order in blockchain").Result()
	}

	if !bytes.Equal(order.Sender, msg.Sender) {
		return sdk.NewError(StoreKey, CodeNotMatchSender, "The cancel addr is not match order sender").Result()
	}

	return sdk.Result{}
}

func handleMsgCancelTradingPair(ctx sdk.Context, msg MsgCancelTradingPair, keeper Keeper) sdk.Result {

	if err := checkMsgCancelTradingPair(keeper, msg, ctx); err != nil {
		return err.Result()
	}

	// Add del request to store
	dlk := NewDelistKeeper(keeper.marketKey)
	dlk.AddDelistRequest(ctx, msg.EffectiveTime, msg.Symbol)

	// send msg to kafka
	values := strings.Split(msg.Symbol, SymbolSeparator)
	msgInfo := CancelMarketInfo{
		Stock:   values[0],
		Money:   values[1],
		Deleter: msg.Sender.String(),
		DelTime: msg.EffectiveTime,
	}
	keeper.SendMsg(CancelMarketInfoKey, msgInfo)

	return sdk.Result{Tags: sdk.NewTags(
		Sender, msg.Sender.String(),
		TradingPair, msg.Symbol,
		EffectiveTime, strconv.Itoa(int(msg.EffectiveTime)),
	)}
}

func checkMsgCancelTradingPair(keeper Keeper, msg MsgCancelTradingPair, ctx sdk.Context) sdk.Error {

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	marketParams := keeper.GetParams(ctx)
	currTime := ctx.BlockHeader().Time.Unix()
	if msg.EffectiveTime < currTime+marketParams.MarketMinExpiredTime {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidTime, "Invalid Cancel Time")
	}

	info, err := keeper.GetMarketInfo(ctx, msg.Symbol)
	if err != nil {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, err.Error())
	}

	stockToken := keeper.axk.GetToken(ctx, info.Stock)
	if !bytes.Equal(msg.Sender, stockToken.GetOwner()) {
		return sdk.NewError(CodeSpaceMarket, CodeNotMatchSender, "Not match market info sender")
	}

	return nil
}

func calculateAmount(price, quantity int64, pricePrecision byte) sdk.Dec {
	actualPrice := sdk.NewDec(price).Quo(sdk.NewDec(int64(math.Pow10(int(pricePrecision)))))
	money := actualPrice.Mul(sdk.NewDec(quantity))
	return money.Add(sdk.NewDec(ExtraFrozenMoney)).Ceil()
}

func marketStoreKey(prefix []byte, params ...string) []byte {
	buf := bytes.NewBuffer(prefix)
	for _, param := range params {
		if _, err := buf.Write([]byte(param)); err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}
