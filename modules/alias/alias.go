package alias

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper = keepers.NewKeeper
	DefaultParams = types.DefaultParams
)

type (
	Keeper         = keepers.Keeper
	MsgAliasUpdate = types.MsgAliasUpdate
)
