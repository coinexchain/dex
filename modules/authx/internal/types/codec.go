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
	cdc.RegisterConcrete(AccountX{}, "authx/AccountX", nil)
}
