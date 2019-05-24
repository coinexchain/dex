package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params Params `json:"param"`
}

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
func InitGenesis(ctx sdk.Context, keeper AccountXKeeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper AccountXKeeper) GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) Validate() error {
	minGasPrice := data.Params.MinGasPrice
	if minGasPrice == 0 {
		return ErrInvalidMinGasPrice("invalid min gas price: 0")
	}
	return nil
}
