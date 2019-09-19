package types

import (
	"os"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

// MsgSetMemoRequired tests
func TestMain(m *testing.M) {
	dex.InitSdkConfig()
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
	addr := sdk.AccAddress([]byte("input1"))
	var msg = NewMsgSend(addr, sdk.AccAddress{}, sdk.NewCoins(), 0)
	if actual := msg.GetSigners(); !reflect.DeepEqual(actual, []sdk.AccAddress{addr}) {
		t.Errorf("Msg.GetSigners() = %v, want %v", actual, []sdk.AccAddress{addr})
	}
}

func TestMsgMultiSend_ValidateBasic(t *testing.T) {
	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("input")))
	addr2 := sdk.AccAddress(crypto.AddressHash([]byte("output")))
	cet123 := sdk.NewCoins(sdk.NewInt64Coin("cet", 123))
	cet111 := sdk.NewCoins(sdk.NewInt64Coin("cet", 111))

	tests := []struct {
		name string
		msg  MsgMultiSend
		want sdk.Error
	}{
		{
			"base_case",
			NewMsgMultiSend([]bank.Input{bank.NewInput(addr1, cet123)}, []bank.Output{bank.NewOutput(addr2, cet123)}),
			nil,
		},
		{
			"no_input",
			NewMsgMultiSend([]bank.Input{}, []bank.Output{bank.NewOutput(addr2, cet123)}),
			ErrNoInputs(),
		},
		{
			"no_output",
			NewMsgMultiSend([]bank.Input{bank.NewInput(addr1, cet123)}, []bank.Output{}),
			ErrNoOutputs(),
		},
		{
			"err_coin",
			NewMsgMultiSend([]bank.Input{bank.NewInput(addr1, cet123)}, []bank.Output{bank.NewOutput(addr2, cet111)}),
			ErrInputOutputMismatch("inputs outputs mismatch"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgMultiSend.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgMultiSend(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("address")))
	cet123 := sdk.NewCoins(sdk.NewInt64Coin("cet", 123))

	msg := NewMsgMultiSend([]bank.Input{bank.NewInput(sdk.AccAddress{}, cet123)}, []bank.Output{bank.NewOutput(addr, cet123)})
	require.Error(t, msg.ValidateBasic())
	msg = NewMsgMultiSend([]bank.Input{bank.NewInput(addr, cet123)}, []bank.Output{bank.NewOutput(sdk.AccAddress{}, cet123)})
	require.Error(t, msg.ValidateBasic())
}
