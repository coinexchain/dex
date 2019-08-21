package incentive

import (
	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

type (
	GenesisState = types.GenesisState
	State        = types.State
	Params       = types.Params
	Plan         = types.Plan
	Keeper       = keepers.Keeper
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
)

var (
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
	DefaultParams       = types.DefaultParams
	NewKeeper           = keepers.NewKeeper
)
