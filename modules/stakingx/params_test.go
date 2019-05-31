package stakingx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, "1000000000000", params.MinSelfDelegation.String())
}