package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type TokenKeeper struct {
	paramSubspace params.Subspace
	// TODO
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (tk TokenKeeper) SetParams(ctx sdk.Context, params Params) {
	//tk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (tk TokenKeeper) GetParams(ctx sdk.Context) (params Params) {
	//tk.paramSubspace.GetParamSet(ctx, &params)
	return
}
