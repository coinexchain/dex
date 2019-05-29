package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx"
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
		Address:      addr,
		MemoRequired: true,
		LockedCoins: []authx.LockedCoin{
			{Coin: dex.NewCetCoin(10), UnlockTime: 10},
		},
		FrozenCoins: dex.NewCetCoins(1000),
	}
	app.accountXKeeper.SetAccountX(ctx, accx)

	state := app.exportGenesisState(ctx)
	require.Equal(t, 1, len(state.Accounts))
	require.Equal(t, sdk.NewInt(int64(1000)), state.Accounts[0].Coins.AmountOf("cet"))
	require.Equal(t, true, state.Accounts[0].MemoRequired)
	require.Equal(t, int64(10), state.Accounts[0].LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(int64(10)), state.Accounts[0].LockedCoins[0].Coin.Amount)
	require.Equal(t, "cet", state.Accounts[0].LockedCoins[0].Coin.Denom)
	require.Equal(t, "1000cet", state.Accounts[0].FrozenCoins.String())

}

func TestExportDefaultAccountXState(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initApp(acc)
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight()})

	state := app.exportGenesisState(ctx)
	require.Equal(t, 1, len(state.Accounts))
	require.Equal(t, sdk.NewInt(int64(1000)), state.Accounts[0].Coins.AmountOf("cet"))
	require.Equal(t, false, state.Accounts[0].MemoRequired)
	require.Nil(t, state.Accounts[0].LockedCoins)
	require.Nil(t, state.Accounts[0].FrozenCoins)

}
