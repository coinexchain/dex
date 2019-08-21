package incentive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper keepers.Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Param)
	err := keeper.SetState(ctx, data.State)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper keepers.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	state := keeper.GetState(ctx)
	return types.NewGenesisState(state, params)
}
