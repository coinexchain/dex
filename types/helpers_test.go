package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestResponseFrom(t *testing.T) {
	rsp := ResponseFrom(sdk.ErrOutOfGas("woo"))
	require.Equal(t, "sdk", rsp.Codespace)
	require.Equal(t, uint32(12), rsp.Code)
}

func TestSafeJsonMarshal(t *testing.T) {
	require.Equal(t, []byte("1"), SafeJsonMarshal(1))
	require.Equal(t, []byte{}, SafeJsonMarshal(TestSafeJsonMarshal))
}
