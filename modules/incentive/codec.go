package incentive

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(State{}, "incentive/state", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	msgCdc.Seal()
}
