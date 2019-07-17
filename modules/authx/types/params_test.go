package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_Equal(t *testing.T) {
	param := DefaultParams()
	param2 := NewParams(sdk.MustNewDecFromStr("20.0"))
	b := Equal(param2)
	require.Equal(t, true, b)
}
