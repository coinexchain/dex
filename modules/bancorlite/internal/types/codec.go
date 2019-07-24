package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	ModuleCdc = codec.New()
)

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBancorInit{}, "bancorlite/MsgBancorInit", nil)
	cdc.RegisterConcrete(MsgBancorTrade{}, "bancorlite/MsgBancorTrade", nil)
	cdc.RegisterConcrete(MsgBancorCancel{}, "bancorlite/MsgBancorCancel", nil)
}
