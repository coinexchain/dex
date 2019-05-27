package market

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
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
	CreateMarketSpendCet = types.NewCetCoin(1000)
}

func NewHandler(k Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgCreateMarketInfo:
			return handlerMsgCreateMarketInfo(ctx, msg, k)
		case MsgCreateGTEOrder:
			return handlerMsgCreateGTEOrder(ctx, msg, k)
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

	if err := keeper.bnk.DeductFeeFromAddressAndCollectFeetoIncentive(msg.Creator, sdk.Coins{CreateMarketSpendCet}); err != nil {
		return ErrSendTokenFailed(err.Error()).Result()
	}
	key := marketStoreKey(marketIdetifierPrefix, info.Stock+SymbolSeparator+info.Money)
	value := msgCdc.MustMarshalBinaryBare(info)
	ctx.KVStore(keeper.marketKey).Set(key, value)

	return sdk.Result{}
}

func checkMsgCreateMarketInfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {
	key := marketStoreKey(marketIdetifierPrefix, msg.Stock+SymbolSeparator+msg.Money)
	store := ctx.KVStore(keeper.marketKey)
	if v := store.Get(key); v == nil {
		return ErrNoExistKeyInStore().Result()
	}

	if !keeper.axk.IsTokenExists(msg.Money) || !keeper.axk.IsTokenExists(msg.Stock) {
		return ErrTokenNoExist().Result()
	}

	if !keeper.axk.IsTokenIssuer(msg.Stock, []byte(msg.Creator)) && !keeper.axk.IsTokenIssuer(msg.Money, []byte(msg.Creator)) {
		return ErrInvalidTokenIssuer().Result()
	}

	if msg.PricePrecision < MinimumTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	if !keeper.bnk.HaveSufficientCoins(msg.Creator, sdk.Coins{CreateMarketSpendCet}) {
		return ErrNoHaveSufficientCoins().Result()
	}

	return sdk.Result{}
}

func handlerMsgCreateGTEOrder(ctx sdk.Context, msg MsgCreateGTEOrder, keeper Keeper) sdk.Result {

	store := ctx.KVStore(keeper.marketKey)
	if store == nil {
		return ErrNoStoreEngine().Result()
	}

	if ret := checkMsgCreateGTEOrder(store, msg, keeper); !ret.IsOK() {
		return ret
	}

	//TODO, bech32 encode need to solve.
	addr, err := bech32.Encode("", msg.Sender)
	if err != nil {
		return ErrInvalidAddress().Result()
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

	if msg.Side != match.BUY && msg.Side != match.SELL {
		return ErrInvalidTradeSide().Result()
	}

	values := strings.Split(msg.Symbol, SymbolSeparator)
	if len(values) != 2 {
		return ErrInvalidSymbol().Result()
	}

	denom = values[0]
	if msg.Side == match.BUY {
		denom = values[1]
	}

	if msg.OrderType != LimitOrder {
		return ErrInvalidOrderType().Result()
	}

	if value = store.Get(marketStoreKey(marketIdetifierPrefix, msg.Symbol)); value == nil {
		return ErrNoExistKeyInStore().Result()
	}

	msgCdc.MustUnmarshalBinaryBare(value, &marketInfo)
	if msg.PricePrecision > marketInfo.PricePrecision {
		return ErrInvalidPricePrecision().Result()
	}

	coin := sdk.NewCoin(denom, calculateAmount(msg.Price, msg.Quantity, msg.PricePrecision).RoundInt())
	if !keeper.bnk.HaveSufficientCoins(msg.Sender, sdk.Coins{coin}) {
		return ErrNoHaveSufficientCoins().Result()
	}

	if keeper.axk.IsTokenFrozen(msg.Sender, denom) {
		return ErrTokenFrozenByIssuer().Result()
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
