package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.NoError(t, state.Validate())

	require.Equal(t, "cet", state.GovData.DepositParams.MinDeposit[0].Denom)
	require.Equal(t, "cet", state.CrisisData.ConstantFee.Denom)
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "504h0m0s", state.StakingData.Params.UnbondingTime.String()) // 21 days
	require.Equal(t, 42, int(state.StakingData.Params.MaxValidators))
	require.Equal(t, sdk.NewDec(20), state.AuthXData.Params.MinGasPriceLimit)
}
