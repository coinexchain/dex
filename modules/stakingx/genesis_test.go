package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/testutil"
)

func TestGenesisState_Validate(t *testing.T) {
	//valid state
	validState := GenesisState{
		Params: DefaultParams(),
	}
	require.Nil(t, validState.Validate())

	//invalidMinSelfDelegation
	invalidMinSelfDelegation := GenesisState{
		Params: Params{
			MinSelfDelegation:    sdk.ZeroInt(),
			NonBondableAddresses: []sdk.AccAddress{},
		},
	}
	require.NotNil(t, invalidMinSelfDelegation.Validate())

	//invalidNonBondedAddresses
	nonBondedAddresses := make([]sdk.AccAddress, 2)
	nonBondedAddresses[0] = testutil.ToAccAddress("myaddr")
	nonBondedAddresses[1] = testutil.ToAccAddress("myaddr")
	invalidNonBondedAddresses := GenesisState{
		Params: Params{
			MinSelfDelegation:    sdk.NewInt(DefaultMinSelfDelegation),
			NonBondableAddresses: nonBondedAddresses,
		},
	}
	require.NotNil(t, invalidNonBondedAddresses.Validate())
}

func TestDefaultGenesisState(t *testing.T) {

	defautGenesisState := DefaultGenesisState()
	require.Equal(t, DefaultParams(), defautGenesisState.Params)
}
