package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterConcrete(&BaseToken{}, "asset/BaseToken", nil)
	//cdc.RegisterInterface((*keeper.Keeper)(nil), nil)
	//cdc.RegisterConcrete(&keeper.BaseKeeper{}, "asset/BaseKeeper", nil)
	//cdc.RegisterInterface((*keeper.TokenKeeper)(nil), nil)
	//cdc.RegisterConcrete(&keeper.BaseTokenKeeper{}, "asset/BaseTokenKeeper", nil)

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
	cdc.RegisterConcrete(MsgModifyTokenURL{}, "asset/MsgModifyTokenURL", nil)
	cdc.RegisterConcrete(MsgModifyTokenDescription{}, "asset/MsgModifyTokenDescription", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
