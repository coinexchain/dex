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

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCommentToken{}, "comment/MsgCommentToken", nil)
}
