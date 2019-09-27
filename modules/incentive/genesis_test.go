package incentive_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/incentive"
)

// nolint
func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		State incentive.State
		Param incentive.Params
	}
	field := fields{State: incentive.State{HeightAdjustment: int64(0)}, Param: incentive.DefaultParams()}
	fieldInvalid := fields{State: incentive.State{HeightAdjustment: -1}, Param: incentive.DefaultParams()}

	param1 := incentive.Params{
		DefaultRewardPerBlock: 1,
		Plans: []incentive.Plan{
			{-1, 2, 1, 10}}}

	param2 := incentive.Params{
		DefaultRewardPerBlock: 1,
		Plans: []incentive.Plan{
			{StartHeight: 2, EndHeight: 2, RewardPerBlock: 1, TotalIncentive: 10}}}

	param3 := incentive.Params{
		DefaultRewardPerBlock: 1,
		Plans: []incentive.Plan{
			{2, 20, 0, 10}}}

	param4 := incentive.Params{
		DefaultRewardPerBlock: 1,
		Plans: []incentive.Plan{
			{0, 10, 1, 0}}}

	param5 := incentive.Params{
		DefaultRewardPerBlock: 1,
		Plans: []incentive.Plan{
			{0, 10, 1, 9}}}

	field1 := fields{State: incentive.State{1}, Param: param1}
	field2 := fields{State: incentive.State{1}, Param: param2}
	field3 := fields{State: incentive.State{1}, Param: param3}
	field4 := fields{State: incentive.State{1}, Param: param4}
	field5 := fields{State: incentive.State{1}, Param: param5}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"TestGenesisState_Validate", field, false},
		{"TestGenesisState_Validate_adjustmentHeight", fieldInvalid, true},
		{"TestGenesisState_Invalidate1", field1, true},
		{"TestGenesisState_Invalidate2", field2, true},
		{"TestGenesisState_Invalidate3", field3, true},
		{"TestGenesisState_Invalidate4", field4, true},
		{"TestGenesisState_Invalidate5", field5, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := incentive.GenesisState{
				State:  tt.fields.State,
				Params: tt.fields.Param,
			}
			if err := data.ValidateGenesis(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	state := incentive.DefaultGenesisState()
	require.Equal(t, int64(0), state.State.HeightAdjustment)
	require.Equal(t, int64(10e8), state.Params.Plans[0].RewardPerBlock)
}

func TestExportGenesis(t *testing.T) {
	input := SetupTestInput()
	genesis := incentive.DefaultGenesisState()
	plan := incentive.Plan{EndHeight: 10, RewardPerBlock: 1, TotalIncentive: 10}
	incentive.InitGenesis(input.ctx, input.keeper, genesis)
	err := input.keeper.AddNewPlan(input.ctx, plan)
	require.Nil(t, err)
	gen := incentive.ExportGenesis(input.ctx, input.keeper)
	genesis.Params.Plans = append(genesis.Params.Plans, plan)
	require.Equal(t, genesis, gen)
}
