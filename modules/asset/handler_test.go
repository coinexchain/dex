package asset_test

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
	"reflect"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1E18))
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
			"mint_token_invalid",
			asset.NewMsgMintToken("abc", sdk.NewInt(-1000), owner),
			false,
		},
		{
			"burn_token",
			asset.NewMsgBurnToken("abc", sdk.NewInt(1000), owner),
			true,
		},
		{
			"burn_token_invalid",
			asset.NewMsgBurnToken("abc", sdk.NewInt(-1000), owner),
			false,
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
			"add_token_whitelist_invalid",
			asset.NewMsgAddTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"remove_token_whitelist",
			asset.NewMsgRemoveTokenWhitelist("abc", owner, mockAddrList()),
			true,
		},
		{
			"remove_token_whitelist_invalid",
			asset.NewMsgRemoveTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"forbid_address",
			asset.NewMsgForbidAddr("abc", owner, mockAddrList()),
			true,
		},
		{
			"forbid_address_invalid",
			asset.NewMsgForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"unforbid_address",
			asset.NewMsgUnForbidAddr("abc", owner, mockAddrList()),
			true,
		},
		{
			"unforbid_address_invalid",
			asset.NewMsgUnForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"modify_token_info",
			asset.NewMsgModifyTokenInfo("abc", "www.abc.com", "abc example description", owner),
			true,
		},
		{
			"modify_token_url_invalid",
			asset.NewMsgModifyTokenInfo("abc", string(make([]byte, types.MaxTokenURLLength+1)), "abc example description", owner),
			false,
		},
		{
			"modify_token_description_invalid",
			asset.NewMsgModifyTokenInfo("abc", "www.abc.com", string(make([]byte, types.MaxTokenDescriptionLength+1)), owner),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h(input.ctx, tt.msg); !reflect.DeepEqual(got.IsOK(), tt.want) {
				t.Errorf("handleMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IssueToken_DeductFee(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := asset.NewHandler(input.tk)

	// invalid account issue token
	msg := asset.NewMsgIssueToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", types.TestIdentityString)
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	// issue token deduct fee
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1E18))
	require.NoError(t, err)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	coins := input.tk.GetAccTotalToken(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2100), coins.AmountOf(symbol))
	require.Equal(t, sdk.NewInt(1E18-1E12), coins.AmountOf("cet"))

}

func Test_BurnToken_SubtractCoins(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := asset.NewHandler(input.tk)

	// issue token
	msgIssue := asset.NewMsgIssueToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		true, true, false, false, "", "", types.TestIdentityString)
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1E18))
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
	err := input.tk.AddToken(input.ctx, testAddr, dex.NewCetCoins(1E18))
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
