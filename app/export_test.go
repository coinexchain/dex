package app

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/stretchr/testify/require"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func TestExportGenesisState(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initApp(acc)
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight()})

	accx := authx.AccountX{
		Address: addr, Activated: true, TransferMemoRequired: true}
	app.accountXKeeper.SetAccountX(ctx, accx)

	state := app.exportGenesisState(ctx)
	require.Equal(t, 8, len(state.Accounts))
	require.Equal(t, 8, len(state.AccountsX))
	//require.Equal(t, true, state.AccountsX[0].TransferMemoRequired)
}
