package bancorlite

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper       = keepers.NewKeeper
	NewBancorInfoKeeper = keepers.NewBancorInfoKeeper
	DefaultParams       = types.DefaultParams
)

type (
	Keeper                     = keepers.Keeper
	BancorInfo                 = keepers.BancorInfo
	MsgBancorTradeInfoForKafka = types.MsgBancorTradeInfoForKafka
	MsgBancorInfoForKafka      = types.MsgBancorInfoForKafka
	MsgBancorInit              = types.MsgBancorInit
	MsgBancorTrade             = types.MsgBancorTrade
	MsgBancorCancel            = types.MsgBancorCancel
)
