package incentive

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncentiveCoinsAddress(t *testing.T) {
	require.Equal(t, "cosmos1gc5t98jap4zyhmhmyq5af5s7pyv57w56wjg69u", IncentiveCoinsAccAddr.String())
}
