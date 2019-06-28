package stakingx

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params Params `json:"params"`
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
func (data GenesisState) Validate() error {
	msd := data.Params.MinSelfDelegation
	if !msd.IsPositive() {
		return ErrInvalidMinSelfDelegation(msd)
	}

	addrSet := make(map[string]interface{})
	for _, addr := range data.Params.NonBondableAddresses {
		if _, exists := addrSet[addr.String()]; exists {
			return errors.New("duplicate addresses in non_bondable_addresses")
		}

		addrSet[addr.String()] = nil
	}

	return nil
}
