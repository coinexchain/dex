package market

import (
	"bytes"
	"github.com/btcsuite/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

const (
	MinimumTokenPricePrecision              = 0
	CollectCreateMarketFeeAddress           = ""
	Buy                                     = 1
	Sell                                    = 2
	LimitOrder                    OrderType = 2
)

type OrderType = byte

var CreateMarketSpendCet sdk.Coin

func init() {
	CreateMarketSpendCet = sdk.NewCoin("cet", sdk.NewInt(1000))
}

func NewHandler(k Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgCreateMarketInfo:
			return handlerMsgCreateMarketinfo(ctx, msg, k)
		case MsgCreateGTEOrder:
			return handlerMsgCreateGTEOrder(ctx, msg, k)
		default:
			errMsg := "Unrecognized market Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handlerMsgCreateMarketinfo:
func handlerMsgCreateMarketinfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {

	if ret := checkMsgCreateMarketInfo(msg, keeper); !ret.IsOK() {
		return ret
	}

	info := MarketInfo{
		Stock:             msg.Stock,
		Money:             msg.Money,
		Create:            msg.Creator,
		PricePrecision:    msg.PricePrecision,
		LastExecutedPrice: sdk.NewDec(0),
	}

	key := marketStoreKey(marketIdetifierPrefix, info.Stock+"/"+info.Money)
	value := msgCdc.MustMarshalBinaryBare(info)
	if store := ctx.KVStore(keeper.markeyKey); store != nil {
		store.Set(key, value)
	} else {
		return ErrNoStoreEngine()
	}

	return sdk.Result{}
}

func checkMsgCreateMarketInfo(msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {
	var err error
	if keeper.axk.Exists(msg.Money) != nil || keeper.axk.Exists(msg.Stock) != nil {
		return ErrTokenNoExist()
	}

	if keeper.axk.IsTokenIssuer(msg.Stock, []byte(msg.Creator)) != nil && keeper.axk.IsTokenIssuer(msg.Money, []byte(msg.Creator)) != nil {
		return ErrInvalidTokenIssuer()
	}

	if msg.PricePrecision < MinimumTokenPricePrecision || msg.PricePrecision > sdk.Precision {
		return ErrInvalidPricePrecision()
	}

	//TODO, the deduct fee logic need to discuss.
	if err = keeper.bnk.SendCoins([]byte(msg.Creator), []byte(CollectCreateMarketFeeAddress), []sdk.Coin{CreateMarketSpendCet}); err != nil {
		return ErrSendTokenFailed(err.Error())
	}

	return sdk.Result{}
}

func handlerMsgCreateGTEOrder(ctx sdk.Context, msg MsgCreateGTEOrder, keeper Keeper) sdk.Result {

	store := ctx.KVStore(keeper.markeyKey)
	if store == nil {
		return ErrNoStoreEngine()
	}

	if ret := checkMsgCreateGTEOrder(store, msg, keeper); !ret.IsOK() {
		return ret
	}

	addr, err := bech32.Encode("", msg.Sender)
	if err != nil {
		return ErrInvalidAddress()
	}

	order := Order{
		Sender:         msg.Sender,
		Sequence:       msg.Sequence,
		Symbol:         msg.Symbol,
		OrderType:      msg.OrderType,
		PricePrecision: msg.PricePrecision,
		Price:          sdk.NewDec(msg.Price),
		Quantity:       sdk.NewDec(msg.Quantity),
		Side:           msg.Side,
		TimeInForce:    msg.TimeInForce,
	}
	key := marketStoreKey(orderBookIdetifierPrefix, msg.Symbol, addr+"-"+strconv.Itoa(int(msg.Sequence)))
	value := msgCdc.MustMarshalBinaryBare(order)
	store.Set(key, value)

	return sdk.Result{}
}

func checkMsgCreateGTEOrder(store sdk.KVStore, msg MsgCreateGTEOrder, keeper Keeper) sdk.Result {

	if msg.Side != Buy && msg.Side != Sell {
		return ErrInvalidTradeSide()
	}

	if msg.OrderType != LimitOrder {
		return ErrInvalidOrderType()
	}

	if value := store.Get(marketStoreKey(marketIdetifierPrefix, msg.Symbol)); value == nil {
		return ErrNoExistKeyInStore()
	}
	//TODO. Add additional check condition

	//TODO. Need recompute trader coin number
	coin := sdk.Coins{}
	if !keeper.bnk.HaveSufficientCoins(msg.Sender, coin) {
		return ErrNoHaveSufficientCoins()
	}

	return sdk.Result{}
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
