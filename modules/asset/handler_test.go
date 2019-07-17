package asset

import (
	"reflect"
	"strings"
	"testing"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	input := createTestInput()
	h := NewHandler(input.tk)

	res := h(input.ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized asset Msg type: "))
}

func Test_handleMsg(t *testing.T) {
	input := createTestInput()
	h := NewHandler(input.tk)
	owner, _ := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	err := input.tk.AddToken(input.ctx, testAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  sdk.Msg
		want bool
	}{
		{
			"issue_token",
			NewMsgIssueToken("ABC Token", "abc", 210000000000, testAddr,
				true, true, true, true, "", ""),
			true,
		},
		{
			"issue_token_invalid",
			NewMsgIssueToken("999 Token", "999", 210000000000, testAddr,
				true, true, true, true, "", ""),
			false,
		},
		{
			"transfer_ownership",
			NewMsgTransferOwnership("abc", testAddr, owner),
			true,
		},
		{
			"transfer_ownership_invalid",
			NewMsgTransferOwnership("abc", testAddr, owner),
			false,
		},
		{
			"mint_token",
			NewMsgMintToken("abc", 1000, owner),
			true,
		},
		{
			"mint_token_invalid",
			NewMsgMintToken("abc", -1000, owner),
			false,
		},
		{
			"burn_token",
			NewMsgBurnToken("abc", 1000, owner),
			true,
		},
		{
			"burn_token_invalid",
			NewMsgBurnToken("abc", 9E18+1000, owner),
			false,
		},
		{
			"forbid_token",
			NewMsgForbidToken("abc", owner),
			true,
		},
		{
			"forbid_token_invalid",
			NewMsgForbidToken("abc", testAddr),
			false,
		},
		{
			"unforbid_token",
			NewMsgUnForbidToken("abc", owner),
			true,
		},
		{
			"unforbid_token_invalid",
			NewMsgUnForbidToken("abc", testAddr),
			false,
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", owner, mockAddrList()),
			true,
		},
		{
			"add_token_whitelist_invalid",
			NewMsgAddTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"remove_token_whitelist",
			NewMsgRemoveTokenWhitelist("abc", owner, mockAddrList()),
			true,
		},
		{
			"remove_token_whitelist_invalid",
			NewMsgRemoveTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"forbid_address",
			NewMsgForbidAddr("abc", owner, mockAddrList()),
			true,
		},
		{
			"forbid_address_invalid",
			NewMsgForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"unforbid_address",
			NewMsgUnForbidAddr("abc", owner, mockAddrList()),
			true,
		},
		{
			"unforbid_address_invalid",
			NewMsgUnForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"modify_token_url",
			NewMsgModifyTokenURL("abc", "www.abc.com", owner),
			true,
		},
		{
			"modify_token_url_invalid",
			NewMsgModifyTokenURL("abc", string(make([]byte, 100+1)), owner),
			false,
		},
		{
			"modify_token_description",
			NewMsgModifyTokenDescription("abc", "abc example description", owner),
			true,
		},
		{
			"modify_token_description_invalid",
			NewMsgModifyTokenDescription("abc", string(make([]byte, 1024+1)), owner),
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
	h := NewHandler(input.tk)

	// invalid account issue token
	msg := NewMsgIssueToken("ABC Token", symbol, 210000000000, testAddr,
		false, false, false, false, "", "")
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	// issue token deduct fee
	err := input.tk.AddToken(input.ctx, testAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	coins := input.tk.GetTotalCoins(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(210000000000), coins.AmountOf(symbol))
	require.Equal(t, sdk.NewInt(1E18-1E12), coins.AmountOf("cet"))

}

func Test_BurnToken_SubtractCoins(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := NewHandler(input.tk)

	// issue token
	msgIssue := NewMsgIssueToken("ABC Token", symbol, 2100, testAddr,
		true, true, false, false, "", "")
	err := input.tk.AddToken(input.ctx, testAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// burn token
	msgBurn := NewMsgBurnToken(symbol, 100, testAddr)
	res = h(input.ctx, msgBurn)
	require.True(t, res.IsOK())

	coins := input.tk.GetTotalCoins(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2000), coins.AmountOf(symbol))
}

func Test_MintToken_AddCoins(t *testing.T) {
	input := createTestInput()
	symbol := "abc"
	h := NewHandler(input.tk)

	// issue token
	msgIssue := NewMsgIssueToken("ABC Token", symbol, 2100, testAddr,
		true, true, false, false, "", "")
	err := input.tk.AddToken(input.ctx, testAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// mint token
	msgMint := NewMsgMintToken(symbol, 100, testAddr)
	res = h(input.ctx, msgMint)
	require.True(t, res.IsOK())

	coins := input.tk.GetTotalCoins(input.ctx, testAddr)
	require.Equal(t, sdk.NewInt(2200), coins.AmountOf(symbol))

}
