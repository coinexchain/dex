package comment

import (
	"github.com/coinexchain/dex/modules/comment/internal/types"
	"github.com/coinexchain/dex/modules/comment/internal/keepers"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper = keepers.NewKeeper
)

type (
	Keeper     = keepers.Keeper
)
