package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.NoError(t, state.Validate())

	require.Equal(t, "cet", state.GovData.DepositParams.MinDeposit[0].Denom)
	require.Equal(t, "cet", state.CrisisData.ConstantFee.Denom)
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "504h0m0s", state.StakingData.Params.UnbondingTime.String()) // 21 days
	require.Equal(t, 42, int(state.StakingData.Params.MaxValidators))
	require.True(t, state.AuthXData.Params.MinGasPriceLimit.IsPositive())
}
