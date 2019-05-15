package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
}
