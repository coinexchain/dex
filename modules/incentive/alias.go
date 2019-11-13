package incentive

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

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
	DefaultGenesisState = types.DefaultGenesisState
	DefaultParams       = types.DefaultParams
	NewKeeper           = keepers.NewKeeper
)
