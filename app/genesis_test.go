package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "cet", state.GovData.DepositParams.MinDeposit[0].Denom)
	require.Equal(t, "cet", state.CrisisData.ConstantFee.Denom)
	require.True(t, state.AuthXData.Params.MinGasPriceLimit.IsPositive())
}
