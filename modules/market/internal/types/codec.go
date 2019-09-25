package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Order{}, "market/Order", nil)
	cdc.RegisterConcrete(MarketInfo{}, "market/TradingPair", nil)
	cdc.RegisterConcrete(MsgCreateTradingPair{}, "market/MsgCreateTradingPair", nil)
	cdc.RegisterConcrete(MsgCreateOrder{}, "market/MsgCreateOrder", nil)
	cdc.RegisterConcrete(MsgCancelOrder{}, "market/MsgCancelOrder", nil)
	cdc.RegisterConcrete(MsgCancelTradingPair{}, "market/MsgCancelTradingPair", nil)
	cdc.RegisterConcrete(MsgModifyPricePrecision{}, "market/MsgModifyPricePrecision", nil)
}
