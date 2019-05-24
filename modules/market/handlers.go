package market

import (
	"bytes"
	"github.com/btcsuite/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
	"strconv"
	"strings"
)

const (
	MinimumTokenPricePrecision           = 8
	MaxTokenPricePrecision               = 12
	Buy                                  = 1
	Sell                                 = 2
	LimitOrder                 OrderType = 2
	SymbolSeparator                      = "/"
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

	if err := keeper.bnk.DeductFeeFromAddressAndCollectFeetoIncentive(msg.Creator, sdk.Coins{CreateMarketSpendCet}); err != nil {
		return ErrSendTokenFailed(err.Error())
	}
	key := marketStoreKey(marketIdetifierPrefix, info.Stock+SymbolSeparator+info.Money)
	value := msgCdc.MustMarshalBinaryBare(info)
	ctx.KVStore(keeper.markeyKey).Set(key, value)

	return sdk.Result{}
}

func checkMsgCreateMarketInfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {
	key := marketStoreKey(marketIdetifierPrefix, msg.Stock+SymbolSeparator+msg.Money)
	store := ctx.KVStore(keeper.markeyKey)
	if v := store.Get(key); v == nil {
		return ErrNoExistKeyInStore()
	}

	if !keeper.axk.IsTokenExists(msg.Money) || !keeper.axk.IsTokenExists(msg.Stock) {
		return ErrTokenNoExist()
	}

	if !keeper.axk.IsTokenIssuer(msg.Stock, []byte(msg.Creator)) && !keeper.axk.IsTokenIssuer(msg.Money, []byte(msg.Creator)) {
		return ErrInvalidTokenIssuer()
	}

	if msg.PricePrecision < MinimumTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision()
	}

	if !keeper.bnk.HaveSufficientCoins(msg.Creator, sdk.Coins{CreateMarketSpendCet}) {
		return ErrNoHaveSufficientCoins()
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

	//TODO, bech32 encode need to solve.
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
		Height:         0,
		LeftStock:      sdk.NewDec(0),
		Freeze:         sdk.NewDec(0),
		DealMoney:      sdk.NewDec(0),
		DealStock:      sdk.NewDec(0),
	}
	key := marketStoreKey(orderBookIdetifierPrefix, msg.Symbol, addr+"-"+strconv.Itoa(int(msg.Sequence)))
	value := msgCdc.MustMarshalBinaryBare(order)
	store.Set(key, value)

	return sdk.Result{}
}

func checkMsgCreateGTEOrder(store sdk.KVStore, msg MsgCreateGTEOrder, keeper Keeper) sdk.Result {

	var (
		value      []byte
		denom      string
		marketInfo MarketInfo
	)

	if msg.Side != Buy && msg.Side != Sell {
		return ErrInvalidTradeSide()
	}

	values := strings.Split(msg.Symbol, SymbolSeparator)
	if len(values) != 2 {
		return ErrInvalidSymbol()
	}

	denom = values[0]
	if msg.Side == Buy {
		denom = values[1]
	}

	if msg.OrderType != LimitOrder {
		return ErrInvalidOrderType()
	}

	if value = store.Get(marketStoreKey(marketIdetifierPrefix, msg.Symbol)); value == nil {
		return ErrNoExistKeyInStore()
	}

	msgCdc.MustUnmarshalBinaryBare(value, &marketInfo)
	if msg.PricePrecision > marketInfo.PricePrecision {
		return ErrInvalidPricePrecision()
	}

	coin := sdk.NewCoin(denom, calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision).RoundInt())
	if !keeper.bnk.HaveSufficientCoins(msg.Sender, sdk.Coins{coin}) {
		return ErrNoHaveSufficientCoins()
	}

	if keeper.axk.IsTokenFrozen(msg.Sender, denom) {
		return ErrTokenFrozenByIssuer()
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
