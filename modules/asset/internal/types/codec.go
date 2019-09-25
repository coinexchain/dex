package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterConcrete(&BaseToken{}, "asset/BaseToken", nil)

	cdc.RegisterConcrete(MsgIssueToken{}, "asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "asset/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "asset/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "asset/MsgBurnToken", nil)
	cdc.RegisterConcrete(MsgForbidToken{}, "asset/MsgForbidToken", nil)
	cdc.RegisterConcrete(MsgUnForbidToken{}, "asset/MsgUnForbidToken", nil)
	cdc.RegisterConcrete(MsgAddTokenWhitelist{}, "asset/MsgAddTokenWhitelist", nil)
	cdc.RegisterConcrete(MsgRemoveTokenWhitelist{}, "asset/MsgRemoveTokenWhitelist", nil)
	cdc.RegisterConcrete(MsgForbidAddr{}, "asset/MsgForbidAddr", nil)
	cdc.RegisterConcrete(MsgUnForbidAddr{}, "asset/MsgUnForbidAddr", nil)
	cdc.RegisterConcrete(MsgModifyTokenInfo{}, "asset/MsgModifyTokenInfo", nil)
}
