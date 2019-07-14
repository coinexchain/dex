package incentive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	State State  `json:"state"`
	Param Params `json:"params"`
}

type State struct {
	HeightAdjustment int64 `json:"height_adjustment"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(state State, param Params) GenesisState {
	return GenesisState{
		State: state,
		Param: param,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(State{int64(0)}, DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParam(ctx, data.Param)
	err := keeper.SetState(ctx, data.State)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParam(ctx)
	state := keeper.GetState(ctx)
	return NewGenesisState(state, params)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) ValidateGenesis() error {

	state := data.State
	if state.HeightAdjustment < 0 {
		return sdk.NewError(CodeSpaceIncentive, CodeInvalidAdjustmentHeight, "invalid adjustment Height")
	}
	param := data.Param
	if param.DefaultRewardPerBlock < 0 {
		return sdk.NewError(CodeSpaceIncentive, CodeInvalidDefaultRewardPerBlock, "invalid default reward per block")
	}

	for _, plan := range param.Plans {
		if plan.StartHeight < 0 || plan.EndHeight < 0 {
			return sdk.NewError(CodeSpaceIncentive, CodeInvalidPlanHeight, "invalid incentive plan height")
		}
		if plan.EndHeight <= plan.StartHeight {
			return sdk.NewError(CodeSpaceIncentive, CodeInvalidPlanHeight, "incentive plan end height should be greater than start height")
		}
		if plan.RewardPerBlock <= 0 {
			return sdk.NewError(CodeSpaceIncentive, CodeInvalidRewardPerBlock, "invalid incentive plan reward per block")
		}
		if plan.TotalIncentive <= 0 {
			return sdk.NewError(CodeSpaceIncentive, CodeInvalidTotalIncentive, "invalid incentive plan total incentive reward")
		}
		if (plan.EndHeight-plan.StartHeight)*plan.RewardPerBlock != plan.TotalIncentive {
			return sdk.NewError(CodeSpaceIncentive, CodeInvalidTotalIncentive, "invalid incentive plan")
		}
	}
	return nil
}
