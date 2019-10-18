package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetMemoRequired{}, "bankx/MsgSetMemoRequired", nil)
	cdc.RegisterConcrete(MsgSend{}, "bankx/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "bankx/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgSupervisedSend{}, "bankx/MsgSupervisedSend", nil)
}
