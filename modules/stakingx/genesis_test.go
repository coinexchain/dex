package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	//valid state
	validState := GenesisState{
		Params: DefaultParams(),
	}
	require.Nil(t, validState.ValidateGenesis())

	//invalidMinSelfDelegation
	invalidMinSelfDelegation := GenesisState{
		Params: Params{
			MinSelfDelegation: sdk.ZeroInt(),
		},
	}
	require.NotNil(t, invalidMinSelfDelegation.ValidateGenesis())
}

func TestDefaultGenesisState(t *testing.T) {
	defaultGenesisState := DefaultGenesisState()
	require.Equal(t, DefaultParams(), defaultGenesisState.Params)
}
