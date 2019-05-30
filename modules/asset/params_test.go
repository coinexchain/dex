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
				types.NewCetCoins(TransferOwnershipFee),
				types.NewCetCoins(ForbidAddrFee),
				types.NewCetCoins(UnForbidAddrFee),
				types.NewCetCoins(ForbidTokenFee),
				types.NewCetCoins(UnForbidTokenFee),
				types.NewCetCoins(TokenForbidWhitelistAddFee),
				types.NewCetCoins(TokenForbidWhitelistRemoveFee),
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
				tt.p.TransferOwnershipFee,
				tt.p.ForbidAddrFee,
				tt.p.UnForbidAddrFee,
				tt.p.ForbidTokenFee,
				tt.p.UnForbidTokenFee,
				tt.p.TokenForbidWhitelistAddFee,
				tt.p.TokenForbidWhitelistRemoveFee,
				tt.p.BurnFee,
				tt.p.MintFee,
			}
			if err := p.ValidateGenesis(); (err != nil) != tt.wantErr {
				t.Errorf("Params.ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
