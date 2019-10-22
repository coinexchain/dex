package asset_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/asset"
)

func TestGenesis(t *testing.T) {
	input := createTestInput()
	owner, _ := sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	input.tk.SetParams(input.ctx, asset.DefaultParams())

	state := asset.DefaultGenesisState()

	cet := &asset.BaseToken{
		Name:             "CoinEx Chain Native Token",
		Symbol:           "cet",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            owner,
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  true,
		TokenForbiddable: true,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
		Identity:         asset.TestIdentityString,
	}
	abc := &asset.BaseToken{
		Name:             "ABC Chain Native Token",
		Symbol:           "abc",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            owner,
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  true,
		TokenForbiddable: true,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
		Identity:         asset.TestIdentityString,
	}
	abcDump := &asset.BaseToken{
		Name:             "ABC Chain Native Token",
		Symbol:           "abc",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            owner,
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  true,
		TokenForbiddable: true,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
		Identity:         asset.TestIdentityString,
	}
	abcInvalid := &asset.BaseToken{
		Name:             "ABC Chain Native Token",
		Symbol:           "933",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            owner,
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  true,
		TokenForbiddable: true,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
		Identity:         asset.TestIdentityString,
	}
	state.Tokens = append(state.Tokens, cet, abc, abcDump, abcInvalid)
	require.Error(t, asset.ValidateGenesis(state))
	state.Tokens = state.Tokens[:2]

	whitelist := []string{"cet:coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke"}
	state.Whitelist = append(state.Whitelist, whitelist...)

	forbiddenList := []string{"abc:coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g"}
	state.ForbiddenAddresses = append(state.ForbiddenAddresses, forbiddenList...)

	require.NoError(t, asset.ValidateGenesis(state))
	asset.InitGenesis(input.ctx, input.tk, state)

	res := input.tk.GetWhitelist(input.ctx, "cet")
	require.Equal(t, 1, len(res))
	require.Equal(t, "coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke", res[0].String())

	res = input.tk.GetForbiddenAddresses(input.ctx, "abc")
	require.Equal(t, 1, len(res))
	require.Equal(t, "coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g", res[0].String())

	export := asset.ExportGenesis(input.ctx, input.tk)
	require.Equal(t, int64(asset.DefaultIssueTokenFee), export.Params.IssueTokenFee)
	require.Equal(t, int64(asset.DefaultIssueRareTokenFee), export.Params.IssueRareTokenFee)
	require.Equal(t, 2, len(export.Tokens))
	require.Equal(t, whitelist, export.Whitelist)
	require.Equal(t, forbiddenList, export.ForbiddenAddresses)

	forbiddenList = []string{"abc:coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"}
	state.ForbiddenAddresses = append(state.ForbiddenAddresses, forbiddenList...)
	require.Error(t, asset.ValidateGenesis(state))
}
