package asset

import (
	"testing"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParams_Equal(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()
	require.Equal(t, p1, p2)

	// mount should equal
	cet := types.NewCetCoins(10)
	p1.IssueTokenFee = cet
	require.NotEqual(t, p1, p2)

	// denom should equal
	abc := NewTokenCoins("abc", 1E12)
	p1.IssueTokenFee = abc
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
				types.NewCetCoins(FreezeAddrFee),
				types.NewCetCoins(UnFreezeAddrFee),
				types.NewCetCoins(FreezeTokenFee),
				types.NewCetCoins(UnFreezeTokenFee),
				types.NewCetCoins(TokenFreezeWhitelistAddFee),
				types.NewCetCoins(TokenFreezeWhitelistRemoveFee),
				types.NewCetCoins(BurnFee),
				types.NewCetCoins(MintFee),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				tt.p.IssueTokenFee,
				tt.p.FreezeAddrFee,
				tt.p.UnFreezeAddrFee,
				tt.p.FreezeTokenFee,
				tt.p.UnFreezeTokenFee,
				tt.p.TokenFreezeWhitelistAddFee,
				tt.p.TokenFreezeWhitelistRemoveFee,
				tt.p.BurnFee,
				tt.p.MintFee,
			}
			if err := p.ValidateGenesis(); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
