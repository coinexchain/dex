package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParams_Equal(t *testing.T) {
	param := DefaultParams()
	param2 := NewParams(sdk.MustNewDecFromStr("20.0"))
	b := param.Equal(param2)
	require.Equal(t, true, b)
}
