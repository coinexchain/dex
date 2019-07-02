package incentive

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		State State
		Param Params
	}
	field := fields{State: State{int64(0)}, Param: DefaultParams()}
	fieldInvalid := fields{State: State{-1}, Param: DefaultParams()}

	param1 := Params{
		1,
		[]Plan{
			{-1, 2, 1, 10}}}

	param2 := Params{
		1,
		[]Plan{
			{2, 2, 1, 10}}}

	param3 := Params{
		1,
		[]Plan{
			{2, 20, 0, 10}}}

	param4 := Params{
		1,
		[]Plan{
			{0, 10, 1, 0}}}

	param5 := Params{
		1,
		[]Plan{
			{0, 10, 1, 9}}}

	field1 := fields{State: State{1}, Param: param1}
	field2 := fields{State: State{1}, Param: param2}
	field3 := fields{State: State{1}, Param: param3}
	field4 := fields{State: State{1}, Param: param4}
	field5 := fields{State: State{1}, Param: param5}

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
			data := GenesisState{
				State: tt.fields.State,
				Param: tt.fields.Param,
			}
			if err := data.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	state := DefaultGenesisState()
	require.Equal(t, int64(0), state.State.HeightAdjustment)
	require.Equal(t, int64(10), state.Param.Plans[0].RewardPerBlock)
}

func TestExportGenesis(t *testing.T) {
	input := SetupTestInput()
	genesis := DefaultGenesisState()
	plan := Plan{0, 10, 1, 10}
	InitGenesis(input.ctx, input.keeper, genesis)
	input.keeper.AddNewPlan(input.ctx, plan)
	gen := ExportGenesis(input.ctx, input.keeper)
	genesis.Param.Plans = append(genesis.Param.Plans, plan)
	require.Equal(t, genesis, gen)
}
