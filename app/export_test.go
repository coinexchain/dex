package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func TestExportRestore(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// export
	app1 := initAppWithBaseAccounts(acc)
	ctx1 := app1.NewContext(false, abci.Header{Height: app1.LastBlockHeight()})
	genState1 := app1.exportGenesisState(ctx1)

	// restore & reexport
	app2 := initApp(func(genState *GenesisState) {
		*genState = genState1
	})
	ctx2 := app2.NewContext(false, abci.Header{Height: app2.LastBlockHeight()})
	genState2 := app2.exportGenesisState(ctx2)

	// check
	json1, err1 := codec.MarshalJSONIndent(app1.cdc, genState1)
	json2, err2 := codec.MarshalJSONIndent(app2.cdc, genState2)
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, json1, json2)
}

func TestExportGenesisState(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	amount := cetToken().GetTotalSupply()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}

	// app
	app := initAppWithBaseAccounts(acc)
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
	require.Equal(t, sdk.NewInt(amount), state.Accounts[0].Coins.AmountOf("cet"))
	require.Equal(t, true, state.Accounts[0].MemoRequired)
	require.Equal(t, int64(10), state.Accounts[0].LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(int64(10)), state.Accounts[0].LockedCoins[0].Coin.Amount)
	require.Equal(t, "cet", state.Accounts[0].LockedCoins[0].Coin.Denom)
	require.Equal(t, "1000cet", state.Accounts[0].FrozenCoins.String())
	require.True(t, state.StakingXData.Params.MinSelfDelegation.IsPositive())
}

func TestExportDefaultAccountXState(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	amount := cetToken().GetTotalSupply()

	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}

	// app
	app := initAppWithBaseAccounts(acc)
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight()})

	state := app.exportGenesisState(ctx)
	require.Equal(t, 1, len(state.Accounts))
	require.Equal(t, sdk.NewInt(amount), state.Accounts[0].Coins.AmountOf("cet"))
	require.Equal(t, false, state.Accounts[0].MemoRequired)
	require.Nil(t, state.Accounts[0].LockedCoins)
	require.Nil(t, state.Accounts[0].FrozenCoins)
}

func TestExportAppStateAndValidators(t *testing.T) {
	sk, pk, addr := testutil.KeyPubAddr()
	amount := cetToken().GetTotalSupply()

	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}

	// init app
	app := initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, acc)
		genState.StakingXData.Params.MinSelfDelegation = sdk.NewInt(1e8)
	})

	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})
	ctx := app.NewContext(false, abci.Header{Height: 1})

	// create validator & self delegate minSelfDelegate CET
	valAddr := sdk.ValAddress(addr)
	minSelfDelegate := app.stakingXKeeper.GetParams(ctx).MinSelfDelegation
	createValMsg := testutil.NewMsgCreateValidatorBuilder(valAddr, pk).
		MinSelfDelegation(minSelfDelegate.Int64()).SelfDelegation(minSelfDelegate.Int64()).
		Build()
	createValTx := newStdTxBuilder().
		Msgs(createValMsg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, sk).Build()
	createValResult := app.Deliver(createValTx)
	require.Equal(t, sdk.CodeOK, createValResult.Code)

	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	//next block
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})
	ctx = app.NewContext(false, abci.Header{Height: app.LastBlockHeight() + 1})

	exportState, valset, err := app.ExportAppStateAndValidators(true, []string{})
	require.Nil(t, err)

	var appState GenesisState
	err = app.cdc.UnmarshalJSON(exportState, &appState)
	require.Nil(t, err)

	val := appState.StakingData.Validators
	require.Equal(t, pk, valset[0].PubKey)
	require.Equal(t, val[0].GetTendermintPower(), valset[0].Power)

	valAcc := app.accountKeeper.GetAccount(ctx, addr)
	require.Equal(t, int64(0), appState.DistrData.FeePool.CommunityPool.AmountOf("cet").Int64())
	require.Equal(t, cetToken().GetTotalSupply()-minSelfDelegate.Int64()-100, valAcc.GetCoins().AmountOf("cet").Int64())

	feesCollected := appState.AuthData.CollectedFees
	require.Equal(t, "100cet", feesCollected.String())
}
