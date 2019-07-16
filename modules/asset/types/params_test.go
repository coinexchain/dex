package types

import (
	"github.com/coinexchain/dex/modules/asset"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/types"
)

func TestParams_Equal(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()
	require.Equal(t, p1, p2)
	require.True(t, Equal(p2))

	// mount should equal
	cet := types.NewCetCoins(10)
	IssueTokenFee = cet
	require.NotEqual(t, p1, p2)

	// denom should equal
	abc := asset.newTokenCoins("abc", 1E12)
	IssueTokenFee = abc
	require.NotEqual(t, p1, p2)
}

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
			if err := ValidateGenesis(); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
