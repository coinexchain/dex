package authx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidate(t *testing.T) {
	genState := DefaultGenesisState()
	require.Nil(t, genState.Validate())

	errGenState := NewGenesisState(NewParams(sdk.NewDec(-1)))
	require.NotNil(t, errGenState.Validate())
}
