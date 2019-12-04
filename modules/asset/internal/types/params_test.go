package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	dex "github.com/coinexchain/dex/types"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

func TestGetIssueTokenFee(t *testing.T) {
	p := DefaultParams()
	require.Equal(t, int64(DefaultIssue2CharTokenFee), p.GetIssueTokenFee("aa"))
	require.Equal(t, int64(DefaultIssue3CharTokenFee), p.GetIssueTokenFee("aaa"))
	require.Equal(t, int64(DefaultIssue4CharTokenFee), p.GetIssueTokenFee("aaaa"))
	require.Equal(t, int64(DefaultIssue5CharTokenFee), p.GetIssueTokenFee("aaaaa"))
	require.Equal(t, int64(DefaultIssue6CharTokenFee), p.GetIssueTokenFee("aaaaaa"))
	require.Equal(t, int64(DefaultIssueLongTokenFee), p.GetIssueTokenFee("aaaaaaa"))
}

func TestParams_Equal(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()
	require.Equal(t, p1, p2)
	require.True(t, p1.Equal(p2))

	// mount should equal
	p1.IssueTokenFee = 10
	require.NotEqual(t, p1, p2)

	// denom should equal
	//abc := NewTokenCoins("abc", sdk.NewInt(1e12))
	//p1.IssueTokenFee = abc
	//require.NotEqual(t, p1, p2)
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
			Params{}, // all zeros
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
