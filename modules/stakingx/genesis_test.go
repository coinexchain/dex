package stakingx_test

import (
	"github.com/coinexchain/dex/modules/stakingx"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	//valid state
	validState := stakingx.GenesisState{
		Params: stakingx.DefaultParams(),
	}
	require.Nil(t, validState.ValidateGenesis())

	//invalidMinSelfDelegation
	invalidMinSelfDelegation := stakingx.GenesisState{
		Params: stakingx.Params{
			MinSelfDelegation: sdk.ZeroInt(),
		},
	}
	require.NotNil(t, invalidMinSelfDelegation.ValidateGenesis())
}

func TestDefaultGenesisState(t *testing.T) {
	defaultGenesisState := stakingx.DefaultGenesisState()
	require.Equal(t, stakingx.DefaultParams(), defaultGenesisState.Params)
}
