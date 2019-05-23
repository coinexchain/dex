package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateMarketInfo{}, "cet-chain/MsgCreateMarketInfo", nil)
	cdc.RegisterConcrete(MsgCreateGTEOrder{}, "cet-chain/MsgCreateGTEOrder", nil)
}
