package market

import (
	"github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
)

var (
	NewBaseKeeper = keepers.NewKeeper
	StoreKey      = types.StoreKey
	ModuleName    = types.ModuleName
)

type (
	Keeper = keepers.Keeper
)
