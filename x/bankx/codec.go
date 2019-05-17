package bankx

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetTransferMemoRequired{}, "cet-chain/MsgSetTransferMemoRequired", nil)
}
