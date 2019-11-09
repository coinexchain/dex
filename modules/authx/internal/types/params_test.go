package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

func TestParams_Equal(t *testing.T) {
	codec.RunInitFuncList()
	param := DefaultParams()
	param2 := NewParams(sdk.MustNewDecFromStr("20.0"))
	b := param.Equal(param2)
	require.Equal(t, true, b)
}
