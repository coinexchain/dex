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
	cdc.RegisterConcrete(MsgFreezeAddress{}, "asset/MsgFreezeAddress", nil)
	cdc.RegisterConcrete(MsgUnfreezeAddress{}, "asset/MsgUnfreezeAddress", nil)
	cdc.RegisterConcrete(MsgFreezeToken{}, "asset/MsgFreezeToken", nil)
	cdc.RegisterConcrete(MsgUnfreezeToken{}, "asset/MsgUnfreezeToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "asset/MsgBurnToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "asset/MsgMintToken", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
