package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	paramSubspace params.Subspace

	sk *staking.Keeper

	dk DistributionKeeper
}

func NewKeeper(paramSubspace params.Subspace, sk *staking.Keeper, dk DistributionKeeper) Keeper {
	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
		sk:            sk,
		dk:            dk,
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
