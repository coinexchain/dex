package bankx

import (
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	activationFee := data.Params.ActivationFee
	if activationFee < 0 {
		return sdk.NewError(types.CodeSpaceBankx, types.CodeInvalidActivationFee, "invalid activated fees")
	}
	if lockCoinsFee := data.Params.LockCoinsFee; lockCoinsFee < 0 {
		return sdk.NewError(types.CodeSpaceBankx, types.CodeInvalidLockCoinsFee, "invalid lock coins fee")
	}
	return nil
}
