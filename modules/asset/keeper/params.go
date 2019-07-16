package keeper

import (
	"github.com/coinexchain/dex/modules/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// SetParams sets the asset module's parameters.
func (keeper BaseKeeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (keeper BaseKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	keeper.paramSubspace.GetParamSet(ctx, &params)
	return
}
