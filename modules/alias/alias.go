package alias

import (
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper = keepers.NewKeeper
	DefaultParams = types.DefaultParams
)

type (
	Keeper = keepers.Keeper
)
