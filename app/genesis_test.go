package app

import (
	"testing"

	"github.com/coinexchain/dex/modules/authx"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "cet", state.GovData.DepositParams.MinDeposit[0].Denom)
	require.Equal(t, "cet", state.CrisisData.ConstantFee.Denom)
	require.Equal(t, authx.DefaultMinGasPrice, state.AuthXData.Params.MinGasPrice)
}
