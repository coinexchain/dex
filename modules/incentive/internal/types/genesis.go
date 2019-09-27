package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	State  State  `json:"state"`
	Params Params `json:"params"`
}

type State struct {
	HeightAdjustment int64 `json:"height_adjustment"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(state State, param Params) GenesisState {
	return GenesisState{
		State:  state,
		Params: param,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(State{int64(0)}, DefaultParams())
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) ValidateGenesis() error {
	state := data.State
	if state.HeightAdjustment < 0 {
		return sdk.NewError(CodeSpaceIncentive, CodeInvalidAdjustmentHeight, "invalid adjustment Height")
	}
	param := data.Params
	if param.DefaultRewardPerBlock < 0 {
		return sdk.NewError(CodeSpaceIncentive, CodeInvalidDefaultRewardPerBlock, "invalid default reward per block")
	}

	return CheckPlans(param.Plans)

}
