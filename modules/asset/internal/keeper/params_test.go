package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParams_ValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		p       Params
		wantErr bool
	}{
		{
			"base-case",
			DefaultParams(),
			false,
		},
		{
			"case-invalidate",
			Params{
				sdk.Coins{},
				sdk.Coins{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.ValidateGenesis(); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
