package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	paramSubspace params.Subspace
}

func NewKeeper(paramSubspace params.Subspace) Keeper {
	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
	}
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}
