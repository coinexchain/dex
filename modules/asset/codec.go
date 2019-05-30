package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterConcrete(&BaseToken{}, "asset/Token", nil)

	cdc.RegisterConcrete(MsgIssueToken{}, "asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "asset/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgForbidAddress{}, "asset/MsgForbidAddress", nil)
	cdc.RegisterConcrete(MsgUnforbidAddress{}, "asset/MsgUnforbidAddress", nil)
	cdc.RegisterConcrete(MsgForbidToken{}, "asset/MsgForbidToken", nil)
	cdc.RegisterConcrete(MsgUnforbidToken{}, "asset/MsgUnforbidToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "asset/MsgBurnToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "asset/MsgMintToken", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
