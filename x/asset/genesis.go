package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	//TODO:
	return nil
}
