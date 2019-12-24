package app

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/testutil"
	dex "github.com/coinexchain/cet-sdk/types"
)

func TestExportRestore(t *testing.T) {
	app1, genState1 := startAppWithOneAccountThenExport()
	app2, genState2 := startAppFromGenesisThenExport(genState1)

	// check
	json1, err1 := codec.MarshalJSONIndent(app1.cdc, genState1)
	json2, err2 := codec.MarshalJSONIndent(app2.cdc, genState2)
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, json1, json2)
}

func startAppFromGenesisThenExport(genState1 GenesisState) (*CetChainApp, GenesisState) {
	// restore & reexport
	app2 := initApp(func(genState *GenesisState) {
		*genState = genState1
	})

	ctx2 := app2.NewContext(false, abci.Header{Height: app2.LastBlockHeight()})
	genState2 := app2.ExportGenesisState(ctx2)
	return app2, genState2
}

func startAppWithOneAccountThenExport() (*CetChainApp, GenesisState) {
	_, _, addr := testutil.KeyPubAddr()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}
	// export
	app1 := initAppWithBaseAccounts(acc)
	ctx1 := app1.NewContext(false, abci.Header{Height: app1.LastBlockHeight()})
	genState1 := app1.ExportGenesisState(ctx1)
	return app1, genState1
}

func TestExportGenesisState(t *testing.T) {
	amount := cetToken().GetTotalSupply().Int64()
	state := startAppWithAccountX(amount)

	account := findAccount(t, state)
	require.Equal(t, sdk.NewInt(amount), account.Coins.AmountOf("cet"))

	require.Equal(t, 1, len(state.AuthXData.AccountXs))
	accountX := state.AuthXData.AccountXs[0]
	require.Equal(t, true, accountX.MemoRequired)
	require.Equal(t, int64(10), accountX.LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(int64(10)), accountX.LockedCoins[0].Coin.Amount)
	require.Equal(t, "cet", accountX.LockedCoins[0].Coin.Denom)
	require.Equal(t, "1000cet", accountX.FrozenCoins.String())
	require.True(t, state.StakingXData.Params.MinSelfDelegation > 0)
}

func findAccount(t *testing.T, state GenesisState) *genaccounts.GenesisAccount {
	sort.Slice(state.Accounts, func(i, j int) bool {
		return state.Accounts[i].ModuleName < state.Accounts[j].ModuleName
	})

	require.Equal(t, 6, len(state.Accounts))
	require.Equal(t, "", state.Accounts[0].ModuleName)
	require.Equal(t, authx.ModuleName, state.Accounts[1].ModuleName)
	require.Equal(t, staking.BondedPoolName, state.Accounts[2].ModuleName)
	require.Equal(t, distribution.ModuleName, state.Accounts[3].ModuleName)
	require.Equal(t, gov.ModuleName, state.Accounts[4].ModuleName)
	require.Equal(t, staking.NotBondedPoolName, state.Accounts[5].ModuleName)

	return &state.Accounts[0]
}

func startAppWithAccountX(amount int64) GenesisState {
	_, _, addr := testutil.KeyPubAddr()
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
	state := app.ExportGenesisState(ctx)
	return state
}

func TestExportDefaultAccountXState(t *testing.T) {
	_, _, addr := testutil.KeyPubAddr()
	amount := cetToken().GetTotalSupply().Int64()

	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}

	// app
	app := initAppWithBaseAccounts(acc)
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight()})

	state := app.ExportGenesisState(ctx)
	sort.Slice(state.Accounts, func(i, j int) bool {
		return state.Accounts[i].ModuleName < state.Accounts[j].ModuleName
	})

	account := findAccount(t, state)
	require.Equal(t, sdk.NewInt(amount), account.Coins.AmountOf("cet"))

	require.Equal(t, 0, len(state.AuthXData.AccountXs))
}

func TestExportAppStateAndValidators(t *testing.T) {
	amount := cetToken().GetTotalSupply().Int64()

	sk, pk, addr := testutil.KeyPubAddr()
	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}

	app := startAppWithOneValidator(acc, addr, pk, sk, t)

	//next block
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight() + 1})
	app.EndBlock(abci.RequestEndBlock{Height: app.LastBlockHeight() + 1})
	app.Commit()

	exportState, valset, err := app.ExportAppStateAndValidators(true, []string{})
	require.Nil(t, err)

	var appState GenesisState
	err = app.cdc.UnmarshalJSON(exportState, &appState)
	require.Nil(t, err)

	val := appState.StakingData.Validators
	require.Equal(t, pk, valset[0].PubKey)
	require.Equal(t, val[0].ConsensusPower(), valset[0].Power)

	//since totalPreviousPower == 0, all collected fee will send to CommunityPool
	valAcc := app.accountKeeper.GetAccount(ctx, addr)
	require.Equal(t, sdk.NewDec(100), appState.DistrData.FeePool.CommunityPool.AmountOf("cet"))
	minSelfDelegate := app.stakingXKeeper.GetParams(ctx).MinSelfDelegation
	require.Equal(t, cetToken().GetTotalSupply().SubRaw(minSelfDelegate).SubRaw(100),
		valAcc.GetCoins().AmountOf("cet"))

	//DistributionAccount including OutStanding rewards and CommunityPool
	feeCollectAccount := getDistributionAccount(&appState)
	require.NotNil(t, feeCollectAccount)
	require.Equal(t, sdk.NewInt(100).Int64(), feeCollectAccount.Coins.AmountOf(dex.DefaultBondDenom).Int64())
}

func getDistributionAccount(gs *GenesisState) *genaccounts.GenesisAccount {
	for _, acc := range gs.Accounts {
		if acc.ModuleName == distribution.ModuleName {
			return &acc
		}
	}

	return nil
}

func startAppWithOneValidator(acc auth.BaseAccount, addr sdk.AccAddress, pk crypto.PubKey, sk crypto.PrivKey, t *testing.T) *CetChainApp {
	app := initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, acc)
		genState.StakingXData.Params.MinSelfDelegation = 1e8
	})

	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	ctx := app.NewContext(false, abci.Header{Height: 1})

	// create validator & self delegate minSelfDelegate CET
	createValTx := prepareCreateValidatorTx(addr, app, ctx, pk, sk)

	createValResult := app.Deliver(createValTx)
	require.Equal(t, sdk.CodeOK, createValResult.Code)

	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	return app
}

func prepareCreateValidatorTx(addr sdk.AccAddress, app *CetChainApp, ctx sdk.Context, pk crypto.PubKey, sk crypto.PrivKey) auth.StdTx {
	valAddr := sdk.ValAddress(addr)
	minSelfDelegate := app.stakingXKeeper.GetParams(ctx).MinSelfDelegation

	createValMsg := testutil.NewMsgCreateValidatorBuilder(valAddr, pk).
		MinSelfDelegation(minSelfDelegate).SelfDelegation(minSelfDelegate).
		Commission("0.1", "0.1", "0.01").
		Build()
	createValTx := newStdTxBuilder().
		Msgs(createValMsg).
		GasAndFee(1000000, 100).AccNumSeqKey(0, 0, sk).Build()

	return createValTx
}

func TestExportValidatorsUpdateRestore(t *testing.T) {
	sk, pk, addr := testutil.KeyPubAddr()
	amount := cetToken().GetTotalSupply().Int64()

	acc := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(amount)}
	app1 := startAppWithOneValidator(acc, addr, pk, sk, t)

	//next block
	app1.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app1.LastBlockHeight() + 1}})
	ctx := app1.NewContext(false, abci.Header{Height: app1.LastBlockHeight() + 1})

	exportState1 := app1.ExportGenesisState(ctx)

	// restore & reexport
	app2 := initAppWithValidators(exportState1)
	ctx2 := app2.NewContext(false, abci.Header{Height: app2.LastBlockHeight()})
	exportState2 := app2.ExportGenesisState(ctx2)

	// check
	json1, err1 := codec.MarshalJSONIndent(app1.cdc, exportState1)
	json2, err2 := codec.MarshalJSONIndent(app2.cdc, exportState2)
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, json1, json2)

}

func initAppWithValidators(gs GenesisState) *CetChainApp {
	app := newApp()

	var validatorUpdates []abci.ValidatorUpdate
	for _, val := range gs.StakingData.Validators {
		validatorUpdates = append(validatorUpdates, val.ABCIValidatorUpdate())
	}

	genStateBytes, _ := app.cdc.MarshalJSON(gs)
	app.InitChain(abci.RequestInitChain{ChainId: testChainID, AppStateBytes: genStateBytes, Validators: validatorUpdates})

	return app
}
