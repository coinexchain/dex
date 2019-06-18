package bankx

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/coinexchain/dex/cmd"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/testutil"
)

// MsgSetMemoRequired tests
func TestMain(m *testing.M) {
	cmd.InitSdkConfig()
	os.Exit(m.Run())
}

func TestSetMemoRequiredRoute(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	require.Equal(t, msg.Route(), "bankx")
	require.Equal(t, msg.Type(), "set_memo_required")
}

func TestSetMemoRequiredValidation(t *testing.T) {
	validAddr := sdk.AccAddress([]byte("addr"))
	var emptyAddr sdk.AccAddress

	testutil.ValidateBasic(t, []testutil.TestCase{
		{Valid: true, Msg: NewMsgSetTransferMemoRequired(validAddr, true)},
		{Valid: true, Msg: NewMsgSetTransferMemoRequired(validAddr, false)},
		{Valid: false, Msg: NewMsgSetTransferMemoRequired(emptyAddr, true)},
		{Valid: false, Msg: NewMsgSetTransferMemoRequired(emptyAddr, false)},
	})
}

func TestSetMemoRequiredGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("addr")))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	sign := msg.GetSignBytes()

	expected := `{"type":"bankx/MsgSetMemoRequired","value":{"address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","required":true}}`
	require.Equal(t, expected, string(sign))
}

func TestSetMemoRequiredGetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	signers := msg.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, addr, signers[0])
}

func TestMsgSendRoute(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10))
	var msg = NewMsgSend(addr1, addr2, coins, 10)

	require.Equal(t, msg.Route(), "bankx")
	require.Equal(t, msg.Type(), "send")
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	cet123 := sdk.NewCoins(sdk.NewInt64Coin("cet", 123))
	cet0 := sdk.NewCoins(sdk.NewInt64Coin("cet", 0))
	cet123eth123 := sdk.NewCoins(sdk.NewInt64Coin("cet", 123), sdk.NewInt64Coin("eth", 123))
	cet123eth0 := sdk.Coins{sdk.NewInt64Coin("cet", 123), sdk.NewInt64Coin("eth", 0)}
	eth123 := sdk.Coins{sdk.NewInt64Coin("eth", 123)}

	var emptyAddr sdk.AccAddress
	time := time.Now().Unix()
	validTime := time + 1000
	invalidTime := int64(-1000)

	cases := []struct {
		valid bool
		tx    MsgSend
	}{
		{true, NewMsgSend(addr1, addr2, cet123, 0)},       // valid send
		{true, NewMsgSend(addr1, addr2, cet123eth123, 0)}, // valid send with multiple coins
		{false, NewMsgSend(addr1, addr2, cet0, 0)},        // non positive coin
		{false, NewMsgSend(addr1, addr2, cet123eth0, 0)},  // non positive coin in multicoins
		{false, NewMsgSend(emptyAddr, addr2, cet123, 0)},  // empty from addr
		{false, NewMsgSend(addr1, emptyAddr, cet123, 0)},  // empty to addr
		{true, NewMsgSend(addr1, addr2, cet123, validTime)},
		{false, NewMsgSend(addr1, addr2, cet123eth123, invalidTime)},
		{true, NewMsgSend(addr1, addr2, eth123, 0)},
		{true, NewMsgSend(addr1, addr2, eth123, validTime)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("input")))
	addr2 := sdk.AccAddress(crypto.AddressHash([]byte("output")))
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10))
	var msg = NewMsgSend(addr1, addr2, coins, 0)
	res := msg.GetSignBytes()

	expected := `{"type":"bankx/MsgSend","value":{"amount":[{"amount":"10","denom":"cet"}],"from_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","to_address":"coinex1urhghdgxshs9lg850mgyyqawj5lal5z460yvr8","unlock_time":"0"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	var msg = NewMsgSend(sdk.AccAddress([]byte("input1")), sdk.AccAddress{}, sdk.NewCoins(), 0)
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431]")
}
