package asset_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func TestInvalidMsg(t *testing.T) {
	input := createTestInput()
	h := asset.NewHandler(input.tk)

	res := h(input.ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized asset Msg type: "))
}

func Test_handleMsg(t *testing.T) {
	input := createTestInput()
	h := asset.NewHandler(input.tk)
	owner, _ := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1e18))
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  sdk.Msg
		want bool
	}{
		{
			"issue_token",
			asset.NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
				true, true, true, true, "", "", types.TestIdentityString),
			true,
		},
		{
			"issue_token_invalid",
			asset.NewMsgIssueToken("999 Token", "999", sdk.NewInt(2100), testAddr,
				true, true, true, true, "", "", types.TestIdentityString),
			false,
		},
		{
			"transfer_ownership",
			asset.NewMsgTransferOwnership("abc", testAddr, owner),
			true,
		},
		{
			"transfer_ownership_invalid",
			asset.NewMsgTransferOwnership("abc", testAddr, owner),
			false,
		},
		{
			"mint_token",
			asset.NewMsgMintToken("abc", sdk.NewInt(1000), owner),
			true,
		},
		{
			"burn_token",
			asset.NewMsgBurnToken("abc", sdk.NewInt(1000), owner),
			true,
		},
		{
			"forbid_token",
			asset.NewMsgForbidToken("abc", owner),
			true,
		},
		{
			"forbid_token_invalid",
			asset.NewMsgForbidToken("abc", testAddr),
			false,
		},
		{
			"unforbid_token",
			asset.NewMsgUnForbidToken("abc", owner),
			true,
		},
		{
			"unforbid_token_invalid",
			asset.NewMsgUnForbidToken("abc", testAddr),
			false,
		},
		{
			"add_token_whitelist",
			asset.NewMsgAddTokenWhitelist("abc", owner, mockAddrList()),
			true,
		},
		{
			"remove_token_whitelist",
			asset.NewMsgRemoveTokenWhitelist("abc", owner, mockAddrList()),
			true,
		},
		{
			"forbid_address",
			asset.NewMsgForbidAddr("abc", owner, mockAddrListNoOwner()),
			true,
		},
		{
			"unforbid_address",
			asset.NewMsgUnForbidAddr("abc", owner, mockAddrListNoOwner()),
			true,
		},
		{
			"modify_token_info",
			asset.NewMsgModifyTokenInfo("abc", "www.abc.com", "abc example description", types.TestIdentityString, owner,
				types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
				types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
				types.DoNotModifyTokenInfo, types.DoNotModifyTokenInfo, // TODO
			),
			true,
		},
		{
			"modify_token_url_invalid",
			asset.NewMsgModifyTokenInfo("abc", string(make([]byte, types.MaxTokenURLLength+1)), "abc example description", types.TestIdentityString, owner,
				"NewName", "123", "true", "true", "true", "true"),
			false,
		},
		{
			"modify_token_description_invalid",
			asset.NewMsgModifyTokenInfo("abc", "www.abc.com", string(make([]byte, types.MaxTokenDescriptionLength+1)), types.TestIdentityString, owner,
				"NewName", "123", "true", "true", "true", "true"),
			false,
		},
		{
			"modify_token_identity_invalid",
			asset.NewMsgModifyTokenInfo("abc", "www.abc.com", "abc example description", string(make([]byte, types.MaxTokenIdentityLength+1)), owner,
				"NewName", "123", "true", "true", "true", "true"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, h(input.ctx, tt.msg).IsOK())
		})
	}
}

func Test_IssueToken_DeductFee(t *testing.T) {
	testIssueTokenDeductFee(t, "abc")
	testIssueTokenDeductFee(t, "abcd")
	testIssueTokenDeductFee(t, "abcde")
	testIssueTokenDeductFee(t, "abcdef")
	testIssueTokenDeductFee(t, "abcdefg")
	testIssueTokenDeductFee(t, "abcdefgh")
	testIssueTokenDeductFee(t, "abcdefghi")
}

func testIssueTokenDeductFee(t *testing.T, symbol string) {
	input := createTestInput()
	h := asset.NewHandler(input.tk)

	// invalid account issue token
	msg := asset.NewMsgIssueToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	// issue token deduct fee
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1e18))
	require.NoError(t, err)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	coins := input.tk.GetAccTotalToken(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2100), coins.AmountOf(symbol))
	require.Equal(t, sdk.NewInt(1e18-types.DefaultParams().GetIssueTokenFee(symbol)), coins.AmountOf("cet"))
}

func Test_BurnToken_SubtractCoins(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := asset.NewHandler(input.tk)

	// issue token
	msgIssue := asset.NewMsgIssueToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1e18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// burn token
	msgBurn := asset.NewMsgBurnToken(symbol, sdk.NewInt(100), testAddr)
	res = h(input.ctx, msgBurn)
	require.True(t, res.IsOK())

	coins := input.tk.GetAccTotalToken(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2000), coins.AmountOf(symbol))
}

func Test_MintToken_AddCoins(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := asset.NewHandler(input.tk)

	// issue token
	msgIssue := asset.NewMsgIssueToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1e18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// mint token
	msgMint := asset.NewMsgMintToken(symbol, sdk.NewInt(100), testAddr)
	res = h(input.ctx, msgMint)
	require.True(t, res.IsOK())

	coins := input.tk.GetAccTotalToken(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2200), coins.AmountOf(symbol))
}

func Test_CollectTokenModificationInfoOK(t *testing.T) {
	token, err := asset.NewToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
		true, true, true, true, "", "", "id")
	require.NoError(t, err)

	msg := newMsgModifyTokenInfo()
	msg.URL = "new.url"
	newToken := modifyToken(t, token, msg, "")
	require.Equal(t, msg.URL, newToken.GetURL())

	msg = newMsgModifyTokenInfo()
	msg.Description = "newDesc"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, msg.Description, newToken.GetDescription())

	msg = newMsgModifyTokenInfo()
	msg.Identity = "newID"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, msg.Identity, newToken.GetIdentity())

	msg = newMsgModifyTokenInfo()
	msg.Name = "newName"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, msg.Name, newToken.GetName())

	msg = newMsgModifyTokenInfo()
	msg.TotalSupply = "2200"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, sdk.NewInt(2200), newToken.GetTotalSupply())

	msg = newMsgModifyTokenInfo()
	msg.Mintable = "false"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, false, newToken.GetMintable())

	msg = newMsgModifyTokenInfo()
	msg.Burnable = "false"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, false, newToken.GetBurnable())

	msg = newMsgModifyTokenInfo()
	msg.AddrForbiddable = "false"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, false, newToken.GetAddrForbiddable())

	msg = newMsgModifyTokenInfo()
	msg.TokenForbiddable = "false"
	newToken = modifyToken(t, token, msg, "")
	require.Equal(t, false, newToken.GetTokenForbiddable())
}

func Test_CollectTokenModificationInfoErr(t *testing.T) {
	token, err := asset.NewToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
		true, true, true, true, "", "", "id")
	require.NoError(t, err)

	msg := newMsgModifyTokenInfo()
	msg.TotalSupply = "2200a"
	modifyToken(t, token, msg, "invalid token TotalSupply: 2200a")

	msg = newMsgModifyTokenInfo()
	msg.Mintable = "fa1se"
	modifyToken(t, token, msg, "invalid token Mintable: fa1se")

	msg = newMsgModifyTokenInfo()
	msg.Burnable = "truu"
	modifyToken(t, token, msg, "invalid token Burnable: truu")

	msg = newMsgModifyTokenInfo()
	msg.AddrForbiddable = "hello"
	modifyToken(t, token, msg, "invalid token AddrForbiddable: hello")

	msg = newMsgModifyTokenInfo()
	msg.TokenForbiddable = "123"
	modifyToken(t, token, msg, "invalid token TokenForbiddable: 123")
}

func modifyToken(t *testing.T, token types.Token, msg asset.MsgModifyTokenInfo, errMsg string) types.Token {
	newURL, newDesc, newID, newName, newSupply,
		newMintable, newBurnable, newAddrForbiddable, newTokenForbiddable,
		err := asset.CollectTokenModificationInfo(token, msg)
	if errMsg != "" {
		require.Error(t, err)
		require.Contains(t, err.Error(), errMsg)
		return nil
	}

	require.NoError(t, err)
	newToken, err := types.NewToken(newName, token.GetSymbol(), newSupply, token.GetOwner(),
		newMintable, newBurnable, newAddrForbiddable, newTokenForbiddable,
		newURL, newDesc, newID)
	require.NoError(t, err)
	return newToken
}

func newMsgModifyTokenInfo() asset.MsgModifyTokenInfo {
	return asset.MsgModifyTokenInfo{
		URL:              types.DoNotModifyTokenInfo,
		Description:      types.DoNotModifyTokenInfo,
		Identity:         types.DoNotModifyTokenInfo,
		Name:             types.DoNotModifyTokenInfo,
		TotalSupply:      types.DoNotModifyTokenInfo,
		Mintable:         types.DoNotModifyTokenInfo,
		Burnable:         types.DoNotModifyTokenInfo,
		AddrForbiddable:  types.DoNotModifyTokenInfo,
		TokenForbiddable: types.DoNotModifyTokenInfo,
	}
}
