package authx_test

import (
	"github.com/coinexchain/dex/modules/authx"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidate(t *testing.T) {
	genState := authx.DefaultGenesisState()
	require.Nil(t, genState.ValidateGenesis())

	errGenState := authx.NewGenesisState(authx.NewParams(sdk.NewDec(-1)), []authx.AccountX{})
	require.NotNil(t, errGenState.ValidateGenesis())
}

func TestExport(t *testing.T) {
	accx := authx.NewAccountX(sdk.AccAddress([]byte("addr")), false, nil, nil)

	testInput := setupTestInput()
	genState1 := authx.NewGenesisState(authx.NewParams(sdk.NewDec(50)), []authx.AccountX{accx})
	authx.InitGenesis(testInput.ctx, testInput.axk, genState1)
	genState2 := authx.ExportGenesis(testInput.ctx, testInput.axk)
	require.Equal(t, genState1, genState2)
	require.True(t, genState2.Params.Equal(genState1.Params))
}
