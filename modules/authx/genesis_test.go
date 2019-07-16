package authx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidate(t *testing.T) {
	genState := DefaultGenesisState()
	require.Nil(t, genState.ValidateGenesis())

	errGenState := NewGenesisState(NewParams(sdk.NewDec(-1)))
	require.NotNil(t, errGenState.ValidateGenesis())
}

func TestExport(t *testing.T) {
	testInput := setupTestInput()
	genState1 := NewGenesisState(NewParams(sdk.NewDec(50)))
	InitGenesis(testInput.ctx, testInput.axk, genState1)
	genState2 := ExportGenesis(testInput.ctx, testInput.axk)
	require.Equal(t, genState1, genState2)
	require.True(t, genState2.Params.Equal(genState1.Params))
}
