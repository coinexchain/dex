package bancorlite

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
)

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	codec.RunInitFuncList()
	type args struct {
		data json.RawMessage
	}
	tests := []struct {
		name    string
		a       AppModuleBasic
		args    args
		wantErr bool
	}{
		{
			name: "positive",
			a:    AppModuleBasic{},
			args: args{
				data: AppModuleBasic{}.DefaultGenesis(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AppModuleBasic{}
			if err := a.ValidateGenesis(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AppModuleBasic.ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
