package incentive

import "testing"

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		State State
		Param Params
	}
	field := fields{State: State{int64(0)}, Param: DefaultParams()}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "TestGenesisState_Validate", fields: field, wantErr: false},
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
