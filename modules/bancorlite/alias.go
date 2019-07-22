package bancorlite

import (
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper       = keepers.NewKeeper
	NewBancorInfoKeeper = keepers.NewBancorInfoKeeper
)

type (
	Keeper = keepers.Keeper
)
