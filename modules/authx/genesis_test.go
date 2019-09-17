package authx_test

import (
	"github.com/coinexchain/dex/modules/authx"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidate(t *testing.T) {
	addr1, err := sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	addr2, err := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	require.NoError(t, err)

	genState := authx.DefaultGenesisState()
	require.Nil(t, genState.ValidateGenesis())

	genState = authx.NewGenesisState(authx.NewParams(sdk.NewDec(10)), []authx.AccountX{authx.NewAccountXWithAddress(addr1), authx.NewAccountXWithAddress(addr2)})
	require.Nil(t, genState.ValidateGenesis())

	errGenState := authx.NewGenesisState(authx.NewParams(sdk.NewDec(-1)), []authx.AccountX{})
	require.NotNil(t, errGenState.ValidateGenesis())

	errGenState = authx.NewGenesisState(authx.NewParams(sdk.NewDec(10)), []authx.AccountX{authx.NewAccountXWithAddress(sdk.AccAddress{})})
	require.NotNil(t, errGenState.ValidateGenesis())

	errGenState = authx.NewGenesisState(authx.NewParams(sdk.NewDec(10)), []authx.AccountX{authx.NewAccountXWithAddress(addr1), authx.NewAccountXWithAddress(addr1)})
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
