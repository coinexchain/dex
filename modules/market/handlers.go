package market

import (
	"bytes"
	"math"
	"strings"

	"github.com/coinexchain/dex/modules/market/match"
	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MinimumTokenPricePrecision           = 8
	MaxTokenPricePrecision               = 18
	LimitOrder                 OrderType = 2
	SymbolSeparator                      = "/"
)

type OrderType = byte

var CreateMarketSpendCet sdk.Coin

func init() {
	CreateMarketSpendCet = types.NewCetCoin(CreateMarketFee)
}

func NewHandler(k Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgCreateMarketInfo:
			return handlerMsgCreateMarketInfo(ctx, msg, k)
		case MsgCreateOrder:
			return handlerMsgCreateOrder(ctx, msg, k)
		case MsgCancelOrder:
			return handlerMsgCancelOrder(ctx, msg, k)
		default:
			errMsg := "Unrecognized market Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handlerMsgCreateMarketinfo:
func handlerMsgCreateMarketInfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {

	if ret := checkMsgCreateMarketInfo(ctx, msg, keeper); !ret.IsOK() {
		return ret
	}

	info := MarketInfo{
		Stock:             msg.Stock,
		Money:             msg.Money,
		Creator:           msg.Creator,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: sdk.NewDec(0),
	}

	key := marketStoreKey(MarketIdentifierPrefix, info.Stock+SymbolSeparator+info.Money)
	value := keeper.cdc.MustMarshalBinaryBare(info)
	ctx.KVStore(keeper.marketKey).Set(key, value)

	return sdk.Result{Tags: info.GetTags()}
}

func checkMsgCreateMarketInfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {
	key := marketStoreKey(MarketIdentifierPrefix, msg.Stock+SymbolSeparator+msg.Money)
	store := ctx.KVStore(keeper.marketKey)
	if v := store.Get(key); v != nil {
		return ErrInvalidSymbol().Result()
	}

	if !keeper.axk.IsTokenExists(ctx, msg.Money) || !keeper.axk.IsTokenExists(ctx, msg.Stock) {
		return ErrTokenNoExist().Result()
	}

	if !keeper.axk.IsTokenIssuer(ctx, msg.Stock, []byte(msg.Creator)) && !keeper.axk.IsTokenIssuer(ctx, msg.Money, []byte(msg.Creator)) {
		return ErrInvalidTokenIssuer().Result()
	}

	if msg.PricePrecision < MinimumTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	if !keeper.bnk.HasCoins(ctx, msg.Creator, sdk.Coins{CreateMarketSpendCet}) {
		return ErrInsufficientCoins().Result()
	}

	return sdk.Result{}
}

func handlerMsgCreateOrder(ctx sdk.Context, msg MsgCreateOrder, keeper Keeper) sdk.Result {

	store := ctx.KVStore(keeper.marketKey)
	if store == nil {
		return ErrNoStoreEngine().Result()
	}

	if ret := checkMsgCreateOrder(ctx, store, msg, keeper); !ret.IsOK() {
		return ret
	}

	order := Order{
		Sender:      msg.Sender,
		Sequence:    msg.Sequence,
		Symbol:      msg.Symbol,
		OrderType:   msg.OrderType,
		Price:       sdk.NewDec(msg.Price),
		Quantity:    msg.Quantity,
		Side:        msg.Side,
		TimeInForce: msg.TimeInForce,
		Height:      ctx.BlockHeight(),
		LeftStock:   0,
		Freeze:      0,
		DealMoney:   0,
		DealStock:   0,
	}

	ork := NewOrderKeeper(keeper.marketKey, order.Symbol, keeper.cdc)
	if err := ork.Add(ctx, &order); err != nil {
		return err.Result()
	}

	return sdk.Result{Tags: order.GetTagsInOrderCreate()}
}

func checkMsgCreateOrder(ctx sdk.Context, store sdk.KVStore, msg MsgCreateOrder, keeper Keeper) sdk.Result {

	var (
		denom      string
		marketInfo MarketInfo
	)

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	values := strings.Split(msg.Symbol, SymbolSeparator)
	denom = values[0]
	if msg.Side == match.BUY {
		denom = values[1]
	}

	marketInfo, err := keeper.GetMarketInfo(ctx, msg.Symbol)
	if err != nil || msg.PricePrecision > marketInfo.PricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	coin := sdk.NewCoin(denom, calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision).RoundInt())
	if !keeper.bnk.HasCoins(ctx, msg.Sender, sdk.Coins{coin}) {
		return ErrInsufficientCoins().Result()
	}

	if keeper.axk.IsTokenFrozen(ctx, denom) {
		return ErrTokenFrozenByIssuer().Result()
	}

	return sdk.Result{}
}

func handlerMsgCancelOrder(ctx sdk.Context, msg MsgCancelOrder, keeper Keeper) sdk.Result {

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	globalKeeper := NewGlobalOrderKeeper(keeper.marketKey, keeper.cdc)
	order := globalKeeper.QueryOrder(ctx, msg.OrderID)
	if order == nil {
		return sdk.NewError(StoreKey, CodeNotFindOrder, "Not find order in blockchain").Result()
	}

	if !bytes.Equal(order.Sender, msg.Sender) {
		return sdk.NewError(StoreKey, CodeNotMatchOrderSender, "The cancel addr is not match order sender").Result()
	}

	ork := NewOrderKeeper(keeper.marketKey, order.Symbol, keeper.cdc)
	if err := ork.Remove(ctx, order); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func calculateAmount(price, quantity int64, pricePrecision byte) sdk.Dec {
	actualPrice := sdk.NewDec(price).Quo(sdk.NewDec(int64(math.Pow10(int(pricePrecision)))))
	return actualPrice.Mul(sdk.NewDec(quantity))
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
