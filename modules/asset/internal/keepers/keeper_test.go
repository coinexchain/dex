package keepers_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/asset/internal/types"
)

func TestTokenKeeper_IssueToken(t *testing.T) {
	input := createTestInput()

	type args struct {
		ctx sdk.Context
		msg types.MsgIssueToken
	}
	tests := []struct {
		name string
		args args
		want sdk.Error
	}{
		{
			"base-case",
			args{
				input.ctx,
				types.NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
					false, false, false, false, "", "", types.TestIdentityString),
			},
			nil,
		},
		{
			"case-duplicate",
			args{
				input.ctx,
				types.NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
					false, false, false, false, "", "", types.TestIdentityString),
			},
			types.ErrDuplicateTokenSymbol("abc"),
		},
		{
			"case-invalid",
			args{
				input.ctx,
				types.NewMsgIssueToken("ABC Token", "999", sdk.NewInt(2100), testAddr,
					false, false, false, false, "", "", types.TestIdentityString),
			},
			types.ErrInvalidTokenSymbol("999"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := input.tk.IssueToken(
				tt.args.ctx,
				tt.args.msg.Name,
				tt.args.msg.Symbol,
				tt.args.msg.TotalSupply,
				tt.args.msg.Owner,
				tt.args.msg.Mintable,
				tt.args.msg.Burnable,
				tt.args.msg.AddrForbiddable,
				tt.args.msg.TokenForbiddable,
				tt.args.msg.URL,
				tt.args.msg.Description,
				tt.args.msg.Identity,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenKeeper.IssueToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenKeeper_TokenStore(t *testing.T) {
	input := createTestInput()

	// set token
	token1, err := types.NewToken("ABC token", "abc", sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token1)
	require.NoError(t, err)

	token2, err := types.NewToken("XYZ token", "xyz", sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token2)
	require.NoError(t, err)

	// get all tokens
	tokens := input.tk.GetAllTokens(input.ctx)
	require.Equal(t, 2, len(tokens))
	require.Contains(t, []string{"abc", "xyz"}, tokens[0].GetSymbol())
	require.Contains(t, []string{"abc", "xyz"}, tokens[1].GetSymbol())

	// remove token
	input.tk.RemoveToken(input.ctx, token1)

	// get token
	res := input.tk.GetToken(input.ctx, token1.GetSymbol())
	require.Nil(t, res)

}
func TestTokenKeeper_TokenReserved(t *testing.T) {
	input := createTestInput()
	addr, _ := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	expectErr := types.ErrInvalidIssueOwner()

	// issue btc token failed
	err := input.tk.IssueToken(input.ctx, "BTC token", "btc", sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.Equal(t, expectErr, err)

	// issue abc token success
	err = input.tk.IssueToken(input.ctx, "ABC token", "abc", sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)

	// issue cet token success
	err = input.tk.IssueToken(input.ctx, "CET token", "cet", sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)

	// cet owner issue btc token success
	err = input.tk.IssueToken(input.ctx, "BTC token", "btc", sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)

	// only cet owner can issue reserved token
	err = input.tk.IssueToken(input.ctx, "ETH token", "eth", sdk.NewInt(2100), addr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.Equal(t, expectErr, err)

}

func TestTokenKeeper_TransferOwnership(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	var addr1, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.TransferOwnership(input.ctx, symbol, testAddr, addr1)
	require.NoError(t, err)

	// get token
	token := input.tk.GetToken(input.ctx, symbol)
	require.NotNil(t, token)
	require.Equal(t, addr1.String(), token.GetOwner().String())

	//case2: invalid token
	err = input.tk.TransferOwnership(input.ctx, "xyz", testAddr, addr1)
	require.Error(t, err)

	//case3: invalid original owner
	err = input.tk.TransferOwnership(input.ctx, symbol, testAddr, addr1)
	require.Error(t, err)

	//case4: invalid new owner
	err = input.tk.TransferOwnership(input.ctx, symbol, addr1, sdk.AccAddress{})
	require.Error(t, err)
}

func TestTokenKeeper_MintToken(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	var addr, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.MintToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.NoError(t, err)

	token := input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, sdk.NewInt(3100), token.GetTotalSupply())
	require.Equal(t, sdk.NewInt(1000), token.GetTotalMint())

	err = input.tk.MintToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.NoError(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, sdk.NewInt(4100), token.GetTotalSupply())
	require.Equal(t, sdk.NewInt(2000), token.GetTotalMint())

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: un mintable token
	// set token mintable: false
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.MintToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 3: mint invalid token
	err = input.tk.IssueToken(input.ctx, "ABC token", "xyz", sdk.NewInt(2100), testAddr,
		true, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.MintToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 4: only token owner can mint token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), addr,
		true, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.MintToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

}

func TestTokenKeeper_BurnToken(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	var addr, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.NoError(t, err)

	token := input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, sdk.NewInt(1100), token.GetTotalSupply())
	require.Equal(t, sdk.NewInt(1000), token.GetTotalBurn())

	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.NoError(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, sdk.NewInt(100), token.GetTotalSupply())
	require.Equal(t, sdk.NewInt(2000), token.GetTotalBurn())

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: un burnable token
	// set token burnable: false
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 3: burn invalid token
	err = input.tk.IssueToken(input.ctx, "ABC token", "xyz", sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 4: only token owner can burn token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), addr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(1000))
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 5: token total supply limited to > 0
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(2100))
	require.Error(t, err)
	err = input.tk.BurnToken(input.ctx, symbol, testAddr, sdk.NewInt(2200))
	require.Error(t, err)
}

func TestTokenKeeper_ForbidToken(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	var addr, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.NoError(t, err)

	token := input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, true, token.GetIsForbidden())

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: un forbiddable token
	// set token forbiddable: false
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 3: duplicate forbid token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.NoError(t, err)

	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 4: only token owner can forbid token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), addr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)

}

func TestTokenKeeper_UnForbidToken(t *testing.T) {
	input := createTestInput()
	symbol := "abc"

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.ForbidToken(input.ctx, symbol, testAddr)
	require.NoError(t, err)

	token := input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, true, token.GetIsForbidden())

	err = input.tk.UnForbidToken(input.ctx, symbol, testAddr)
	require.NoError(t, err)

	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, false, token.GetIsForbidden())

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: unforbid token before forbid token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	err = input.tk.UnForbidToken(input.ctx, symbol, testAddr)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)
}

func TestTokenKeeper_AddTokenWhitelist(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	whitelist := mockAddrList()

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	token := input.tk.GetToken(input.ctx, symbol)

	err = input.tk.AddTokenWhitelist(input.ctx, symbol, testAddr, whitelist)
	require.NoError(t, err)
	addresses := input.tk.GetWhitelist(input.ctx, symbol)
	for _, addr := range addresses {
		require.Contains(t, whitelist, addr)
	}
	require.Equal(t, len(whitelist), len(addresses))

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: un forbiddable token
	// set token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.AddTokenWhitelist(input.ctx, symbol, testAddr, whitelist)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)
}

func TestTokenKeeper_RemoveTokenWhitelist(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	whitelist := mockAddrList()

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	token := input.tk.GetToken(input.ctx, symbol)

	err = input.tk.AddTokenWhitelist(input.ctx, symbol, testAddr, whitelist)
	require.NoError(t, err)
	addresses := input.tk.GetWhitelist(input.ctx, symbol)
	for _, addr := range addresses {
		require.Contains(t, whitelist, addr)
	}
	require.Equal(t, len(whitelist), len(addresses))

	err = input.tk.RemoveTokenWhitelist(input.ctx, symbol, testAddr, []sdk.AccAddress{whitelist[0]})
	require.NoError(t, err)
	addresses = input.tk.GetWhitelist(input.ctx, symbol)
	require.Equal(t, len(whitelist)-1, len(addresses))
	require.NotContains(t, addresses, whitelist[0])

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: un-forbiddable token
	// set token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.RemoveTokenWhitelist(input.ctx, symbol, testAddr, whitelist)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)
}

func TestTokenKeeper_ForbidAddress(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	mock := mockAddrList()

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, true, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	token := input.tk.GetToken(input.ctx, symbol)

	err = input.tk.ForbidAddress(input.ctx, symbol, testAddr, mock)
	require.NoError(t, err)
	forbidden := input.tk.GetForbiddenAddresses(input.ctx, symbol)
	for _, addr := range forbidden {
		require.Contains(t, mock, addr)
	}
	require.Equal(t, len(mock), len(forbidden))

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: addr un-forbiddable token
	// set token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.ForbidAddress(input.ctx, symbol, testAddr, mock)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)
}

func TestTokenKeeper_UnForbidAddress(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	mock := mockAddrList()

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, true, true, "", "", types.TestIdentityString)
	require.NoError(t, err)
	token := input.tk.GetToken(input.ctx, symbol)

	err = input.tk.ForbidAddress(input.ctx, symbol, testAddr, mock)
	require.NoError(t, err)
	forbidden := input.tk.GetForbiddenAddresses(input.ctx, symbol)
	for _, addr := range forbidden {
		require.Contains(t, mock, addr)
	}
	require.Equal(t, len(mock), len(forbidden))

	err = input.tk.UnForbidAddress(input.ctx, symbol, testAddr, []sdk.AccAddress{mock[0]})
	require.NoError(t, err)
	forbidden = input.tk.GetForbiddenAddresses(input.ctx, symbol)
	require.Equal(t, len(mock)-1, len(forbidden))
	require.NotContains(t, forbidden, mock[0])

	// remove token
	input.tk.RemoveToken(input.ctx, token)

	//case 2: addr un-forbiddable token
	// set token
	err = input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.UnForbidAddress(input.ctx, symbol, testAddr, mock)
	require.Error(t, err)

	// remove token
	input.tk.RemoveToken(input.ctx, token)
}

func TestTokenKeeper_ModifyTokenInfo(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	var addr, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	url := "www.abc.com"
	description := "token abc is a example token"
	identity := types.TestIdentityString

	//case 1: base-case ok
	// set token
	err := input.tk.IssueToken(input.ctx, "ABC token", symbol, sdk.NewInt(2100), testAddr,
		true, false, false, false, "www.abc.org", "abc example description", types.TestIdentityString)
	require.NoError(t, err)

	err = input.tk.ModifyTokenInfo(input.ctx, symbol, testAddr, url, description, identity)
	require.NoError(t, err)
	token := input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, url, token.GetURL())
	require.Equal(t, description, token.GetDescription())

	//case 2: only token owner can modify token info
	err = input.tk.ModifyTokenInfo(input.ctx, symbol, addr, "www.abc.org", "token abc is a example token", identity)
	require.Error(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, url, token.GetURL())
	require.Equal(t, description, token.GetDescription())

	//case 3: invalid url
	err = input.tk.ModifyTokenInfo(input.ctx, symbol, testAddr, string(make([]byte, types.MaxTokenURLLength+1)), types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo)
	require.Error(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, url, token.GetURL())
	require.Equal(t, description, token.GetDescription())

	//case 4: invalid description
	err = input.tk.ModifyTokenInfo(input.ctx, symbol, testAddr, types.DoNotModifyTokenInfo, string(make([]byte, types.MaxTokenDescriptionLength+1)), types.DoNotModifyTokenInfo)
	require.Error(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, url, token.GetURL())
	require.Equal(t, description, token.GetDescription())

	//case 4: invalid identity
	err = input.tk.ModifyTokenInfo(input.ctx, symbol, testAddr, types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, string(make([]byte, types.MaxTokenIdentityLength+1)))
	require.Error(t, err)
	token = input.tk.GetToken(input.ctx, symbol)
	require.Equal(t, url, token.GetURL())
	require.Equal(t, description, token.GetDescription())

}
