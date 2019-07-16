package asset

import (
	types2 "github.com/coinexchain/dex/modules/asset/types"
	"reflect"
	"strings"
	"testing"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(BaseKeeper{})

	res := h(sdk.Context{}, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized asset Msg type: "))
}

func Test_handleMsg(t *testing.T) {
	input := setupTestInput()
	h := NewHandler(input.tk)
	owner, _ := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")

	err := input.tk.AddToken(input.ctx, tAccAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  sdk.Msg
		want bool
	}{
		{
			"issue_token",
			types2.NewMsgIssueToken("ABC Token", "abc", 210000000000, tAccAddr,
				true, true, true, true, "", ""),
			true,
		},
		{
			"issue_token_invalid",
			types2.NewMsgIssueToken("999 Token", "999", 210000000000, tAccAddr,
				true, true, true, true, "", ""),
			false,
		},
		{
			"transfer_ownership",
			types2.NewMsgTransferOwnership("abc", tAccAddr, owner),
			true,
		},
		{
			"transfer_ownership_invalid",
			types2.NewMsgTransferOwnership("abc", tAccAddr, owner),
			false,
		},
		{
			"mint_token",
			types2.NewMsgMintToken("abc", 1000, owner),
			true,
		},
		{
			"mint_token_invalid",
			types2.NewMsgMintToken("abc", -1000, owner),
			false,
		},
		{
			"burn_token",
			types2.NewMsgBurnToken("abc", 1000, owner),
			true,
		},
		{
			"burn_token_invalid",
			types2.NewMsgBurnToken("abc", 9E18+1000, owner),
			false,
		},
		{
			"forbid_token",
			types2.NewMsgForbidToken("abc", owner),
			true,
		},
		{
			"forbid_token_invalid",
			types2.NewMsgForbidToken("abc", tAccAddr),
			false,
		},
		{
			"unforbid_token",
			types2.NewMsgUnForbidToken("abc", owner),
			true,
		},
		{
			"unforbid_token_invalid",
			types2.NewMsgUnForbidToken("abc", tAccAddr),
			false,
		},
		{
			"add_token_whitelist",
			types2.NewMsgAddTokenWhitelist("abc", owner, mockWhitelist()),
			true,
		},
		{
			"add_token_whitelist_invalid",
			types2.NewMsgAddTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"remove_token_whitelist",
			types2.NewMsgRemoveTokenWhitelist("abc", owner, mockWhitelist()),
			true,
		},
		{
			"remove_token_whitelist_invalid",
			types2.NewMsgRemoveTokenWhitelist("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"forbid_address",
			types2.NewMsgForbidAddr("abc", owner, mockAddresses()),
			true,
		},
		{
			"forbid_address_invalid",
			types2.NewMsgForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"unforbid_address",
			types2.NewMsgUnForbidAddr("abc", owner, mockAddresses()),
			true,
		},
		{
			"unforbid_address_invalid",
			types2.NewMsgUnForbidAddr("abc", owner, []sdk.AccAddress{}),
			false,
		},
		{
			"modify_token_url",
			types2.NewMsgModifyTokenURL("abc", "www.abc.com", owner),
			true,
		},
		{
			"modify_token_url_invalid",
			types2.NewMsgModifyTokenURL("abc", string(make([]byte, 100+1)), owner),
			false,
		},
		{
			"modify_token_description",
			types2.NewMsgModifyTokenDescription("abc", "abc example description", owner),
			true,
		},
		{
			"modify_token_description_invalid",
			types2.NewMsgModifyTokenDescription("abc", string(make([]byte, 1024+1)), owner),
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
	input := setupTestInput()
	symbol := "abc"
	h := NewHandler(input.tk)

	// invalid account issue token
	msg := types2.NewMsgIssueToken("ABC Token", symbol, 210000000000, tAccAddr,
		false, false, false, false, "", "")
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	// issue token deduct fee
	err := input.tk.AddToken(input.ctx, tAccAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	coins := input.tk.bkx.GetTotalCoins(input.ctx, tAccAddr)
	require.Equal(t, sdk.NewInt(210000000000), coins.AmountOf(symbol))
	require.Equal(t, sdk.NewInt(1E18-1E12), coins.AmountOf("cet"))

}

func Test_BurnToken_SubtractCoins(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	h := NewHandler(input.tk)

	// issue token
	msgIssue := types2.NewMsgIssueToken("ABC Token", symbol, 2100, tAccAddr,
		true, true, false, false, "", "")
	err := input.tk.AddToken(input.ctx, tAccAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// burn token
	msgBurn := types2.NewMsgBurnToken(symbol, 100, tAccAddr)
	res = h(input.ctx, msgBurn)
	require.True(t, res.IsOK())

	coins := input.tk.bkx.GetTotalCoins(input.ctx, tAccAddr)
	require.Equal(t, sdk.NewInt(2000), coins.AmountOf(symbol))
}

func Test_MintToken_AddCoins(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	h := NewHandler(input.tk)

	// issue token
	msgIssue := types2.NewMsgIssueToken("ABC Token", symbol, 2100, tAccAddr,
		true, true, false, false, "", "")
	err := input.tk.AddToken(input.ctx, tAccAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)
	res := h(input.ctx, msgIssue)
	require.True(t, res.IsOK())

	// mint token
	msgMint := types2.NewMsgMintToken(symbol, 100, tAccAddr)
	res = h(input.ctx, msgMint)
	require.True(t, res.IsOK())

	coins := input.tk.bkx.GetTotalCoins(input.ctx, tAccAddr)
	require.Equal(t, sdk.NewInt(2200), coins.AmountOf(symbol))

}
