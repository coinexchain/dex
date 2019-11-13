package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	ModuleCdc *codec.Codec
)

func init() {
	codec.AddInitFunc(func() {
		ModuleCdc = codec.New()
		RegisterCodec(ModuleCdc)
	})
}

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetMemoRequired{}, "bankx/MsgSetMemoRequired", nil)
	cdc.RegisterConcrete(MsgSend{}, "bankx/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "bankx/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgSupervisedSend{}, "bankx/MsgSupervisedSend", nil)
}
