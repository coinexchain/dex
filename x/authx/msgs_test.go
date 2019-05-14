package authx

import (
	"gitlab.com/cetchain/cetchain/testutil"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSetTransferMemoRequired tests

func TestRoute(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	require.Equal(t, msg.Route(), "authx")
	require.Equal(t, msg.Type(), "set_transfer_memo_required")
}

func TestValidation(t *testing.T) {
	validAddr := sdk.AccAddress([]byte("addr"))
	var emptyAddr sdk.AccAddress

	testutil.ValidateBasic(t, []testutil.TestCase{
		{true, NewMsgSetTransferMemoRequired(validAddr, true)},
		{true, NewMsgSetTransferMemoRequired(validAddr, false)},
		{false, NewMsgSetTransferMemoRequired(emptyAddr, true)},
		{false, NewMsgSetTransferMemoRequired(emptyAddr, false)},
	})
}

func TestGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	sign := msg.GetSignBytes()

	expected := `{"type":"cet-chain/MsgSetTransferMemoRequired","value":{"address":"cosmos1v9jxguspv4h2u","required":true}}`
	require.Equal(t, expected, string(sign))
}

func TestGetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	signers := msg.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, addr, signers[0])
}
