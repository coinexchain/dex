package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, "100000000000000", params.MinSelfDelegation.String())
}
