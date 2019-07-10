package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	msgCdc = codec.New()
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Order{}, "market/Order", nil)
	cdc.RegisterConcrete(MarketInfo{}, "market/TradingPair", nil)
	cdc.RegisterConcrete(MsgCreateMarketInfo{}, "market/MsgCreateMarketInfo", nil)
	cdc.RegisterConcrete(MsgCreateOrder{}, "market/MsgCreateOrder", nil)
	cdc.RegisterConcrete(MsgCancelOrder{}, "market/MsgCancelOrder", nil)
	cdc.RegisterConcrete(MsgCancelMarket{}, "market/MsgCancelMarket", nil)
	cdc.RegisterConcrete(QueryMarketInfo{}, "market/QueryMarketInfo", nil)
}

func init() {
	RegisterCodec(msgCdc)
	msgCdc.Seal()
}
