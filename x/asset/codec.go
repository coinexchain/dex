package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "coinex-chain/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "coinex-chain/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgFreezeAddress{}, "coinex-chain/MsgFreezeAddress", nil)
	cdc.RegisterConcrete(MsgUnfreezeAddress{}, "coinex-chain/MsgUnfreezeAddress", nil)
	cdc.RegisterConcrete(MsgFreezeToken{}, "coinex-chain/MsgFreezeToken", nil)
	cdc.RegisterConcrete(MsgUnfreezeToken{}, "coinex-chain/MsgUnfreezeToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "coinex-chain/MsgBurnToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "coinex-chain/MsgMintToken", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
