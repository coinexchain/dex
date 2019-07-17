package comment

import "github.com/cosmos/cosmos-sdk/codec"

var (
	msgCdc = codec.New()
)

func init() {
	RegisterCodec(msgCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCommentToken{}, "comment/MsgCommentToken", nil)
}
