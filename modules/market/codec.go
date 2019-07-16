package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Order{}, "market/Order", nil)
	cdc.RegisterConcrete(MarketInfo{}, "market/TradingPair", nil)
	cdc.RegisterConcrete(MsgCreateTradingPair{}, "market/MsgCreateTradingPair", nil)
	cdc.RegisterConcrete(MsgCreateOrder{}, "market/MsgCreateOrder", nil)
	cdc.RegisterConcrete(MsgCancelOrder{}, "market/MsgCancelOrder", nil)
	cdc.RegisterConcrete(MsgCancelTradingPair{}, "market/MsgCancelTradingPair", nil)
	cdc.RegisterConcrete(QueryMarketInfo{}, "market/QueryMarketInfo", nil)
	cdc.RegisterConcrete(MsgModifyPricePrecision{}, "MsgModifyPricePrecision", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
