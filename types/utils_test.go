package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func TestConcatKeys(t *testing.T) {
	require.Equal(t, []byte("foobar"), ConcatKeys([]byte("foo"), nil, []byte("bar")))
}

func TestErrUnknownRequest(t *testing.T) {
	result := ErrUnknownRequest("bank", bank.MsgSend{})
	require.True(t, strings.Index(result.Log, "Unrecognized bank Msg type: send") > 0)
}

func TestResponseFrom(t *testing.T) {
	rsp := ResponseFrom(sdk.ErrOutOfGas("woo"))
	require.Equal(t, "sdk", rsp.Codespace)
	require.Equal(t, uint32(12), rsp.Code)
}

func TestSafeJsonMarshal(t *testing.T) {
	require.Equal(t, []byte("1"), SafeJSONMarshal(1))
	require.Equal(t, []byte{}, SafeJSONMarshal(TestSafeJsonMarshal))
}
