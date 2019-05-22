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
	cdc.RegisterConcrete(MsgSetMemoRequired{}, "cet-chain/MsgSetMemoRequired", nil)
	cdc.RegisterConcrete(MsgSendWithUnlockTime{}, "cet-chain/MsgSendWithUnlockTime", nil)
}
