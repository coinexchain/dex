package bankx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params types.Params `json:"params"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(param types.Params) GenesisState {
	return GenesisState{
		Params: param,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams())
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
func (data GenesisState) ValidateGenesis() error {
	if activationFee := data.Params.ActivationFee; activationFee < 0 {
		return types.ErrorInvalidActivatingFee()
	}
	if lockCoinsFee := data.Params.LockCoinsFee; lockCoinsFee < 0 {
		return types.ErrorInvalidLockCoinsFee()
	}
	return nil
}
