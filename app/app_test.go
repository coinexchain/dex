package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

type genesisStateCallback func(state *GenesisState)

func newApp() *CetChainApp {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	return NewCetChainApp(logger, db, nil, true, 10000)
}

func initAppWithBaseAccounts(accs ...auth.BaseAccount) *CetChainApp {
	return initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, accs...)
	})
}

func addGenesisAccounts(genState *GenesisState, accs ...auth.BaseAccount) {
	for _, acc := range accs {
		genAcc := NewGenesisAccount(&acc)
		genState.Accounts = append(genState.Accounts, genAcc)
	}
}

func initApp(cb genesisStateCallback) *CetChainApp {
	app := newApp()

	// genesis state
	genState := NewDefaultGenesisState()
	genState.AuthXData.Params.MinGasPriceLimit = sdk.MustNewDecFromStr("0.00000001")
	if cb != nil {
		cb(&genState)
	}

	// init chain
	genStateBytes, _ := app.cdc.MarshalJSON(genState)
	app.InitChain(abci.RequestInitChain{ChainId: "c1", AppStateBytes: genStateBytes})

	return app
}

//func TestMinGasPrice(t *testing.T) {
//	app := newApp()
//	ctx := app.NewContext(true, abci.Header{})
//	minGasPrice := ctx.MinGasPrices()
//	require.False(t, minGasPrice.IsZero())
//}

func TestRouter(t *testing.T) {
	bApp := bam.NewBaseApp(appName, nil, nil, nil)
	app := &CetChainApp{BaseApp: bApp}
	app.registerMessageRoutes()
	require.Nil(t, app.Router().Route("bank"))
}

func TestSend(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10000000000), sdk.NewInt64Coin("eth", 100000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// issue cet token
	ctx := app.NewContext(false, header)
	issueCetmsg := asset.NewMsgIssueToken("CET", "cet", 10000000000000000, sdk.AccAddress("fromaddr"), false, false, false, false)
	app.assetKeeper.IssueToken(ctx, issueCetmsg)

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, time.Now().Unix()+10000)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, bankx.CodeCetCantBeLocked, result.Code)

	msg = bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx = testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result = app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
}

func TestMemo(t *testing.T) {
	key, _, addr := testutil.KeyPubAddr()
	acc0 := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	msgSetMemoRequired := bankx.NewMsgSetTransferMemoRequired(addr, true)
	tx1 := testutil.NewStdTxBuilder("c1").
		Msgs(msgSetMemoRequired).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()
	result1 := app.Deliver(tx1)
	require.Equal(t, errors.CodeOK, result1.Code)

	coins := dex.NewCetCoins(100)
	msgSend := bankx.NewMsgSend(addr, addr, coins, 0)
	tx2 := testutil.NewStdTxBuilder("c1").
		Msgs(msgSend).GasAndFee(1000000, 100).AccNumSeqKey(0, 1, key).Build()

	result2 := app.Deliver(tx2)
	require.Equal(t, bankx.CodeMemoMissing, result2.Code)
}

func TestGasFeeDeductedWhenTxFailed(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// issue cet token
	ctx := app.NewContext(false, header)
	issueCetmsg := asset.NewMsgIssueToken("CET", "cet", 10000000000000000, sdk.AccAddress("fromaddr"), false, false, false, false)
	app.assetKeeper.IssueToken(ctx, issueCetmsg)

	// deliver tx
	coins = dex.NewCetCoins(100000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeInsufficientCoins, result.Code)

	// end block & commit
	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	// check coins
	ctx = app.NewContext(true, abci.Header{})
	require.Equal(t, int64(10000000000-100),
		app.accountKeeper.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet").Int64())
}

func TestSendFromIncentiveAddr(t *testing.T) {
	key, _, fromAddr := testutil.KeyPubAddr()
	incentive.IncentiveCoinsAccAddr = fromAddr
	toAddr := sdk.AccAddress([]byte("addr"))
	coins := dex.NewCetCoinsE8(100)
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeUnauthorized, result.Code)
}

func TestMinSelfDelegation(t *testing.T) {
	key0, pubKey0, addr0 := testutil.KeyPubAddr()
	coins := dex.NewCetCoins(1000)
	acc0 := auth.BaseAccount{Address: addr0, Coins: coins}
	val0 := sdk.ValAddress(addr0)

	// init app
	app := initApp(func(genState *GenesisState) {
		genState.Accounts = append(genState.Accounts, NewGenesisAccountI(&acc0))
		genState.StakingXData.Params.MinSelfDelegation = sdk.NewInt(500)
	})

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver
	msg := testutil.NewMsgCreateValidatorBuilder(val0, pubKey0).
		MinSelfDelegation(400).SelfDelegation(450).
		Build()
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key0).Build()

	result := app.Deliver(tx)
	//require.Nil(t, result.Codespace)
	require.Equal(t, stakingx.CodeMinSelfDelegationBelowRequired, result.Code)
}

func TestDelegatorShares(t *testing.T) {
	// prepare accounts
	valKey, valAcc := testutil.NewBaseAccount(10000, 0, 0)
	valAddr := sdk.ValAddress(valAcc.Address)
	del1Key, del1Acc := testutil.NewBaseAccount(10000, 1, 0)
	del2Key, del2Acc := testutil.NewBaseAccount(10000, 2, 0)

	// init app
	app := initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, valAcc, del1Acc, del2Acc)
		genState.StakingXData.Params.MinSelfDelegation = sdk.NewInt(1)
	})
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// create validator & self delegate 100 CET
	createValMsg := testutil.NewMsgCreateValidatorBuilder(valAddr, valAcc.PubKey).
		MinSelfDelegation(1).SelfDelegation(100).
		Build()
	createValTx := testutil.NewStdTxBuilder("c1").
		Msgs(createValMsg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, valKey).Build()
	createValResult := app.Deliver(createValTx)
	require.Equal(t, sdk.CodeOK, createValResult.Code)

	// delegator1 delegate 100 CET
	del1Msg := staking.NewMsgDelegate(del1Acc.Address, valAddr, dex.NewCetCoin(100))
	del1Tx := testutil.NewStdTxBuilder("c1").
		Msgs(del1Msg).GasAndFee(1000000, 100).AccNumSeqKey(1, 0, del1Key).Build()
	del1Result := app.Deliver(del1Tx)
	require.Equal(t, sdk.CodeOK, del1Result.Code)

	// simulate slash (50 CET)
	ctx := app.NewContext(false, abci.Header{Height: 1})
	val, found := app.stakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, found)
	val.Tokens = val.Tokens.SubRaw(50)
	app.stakingKeeper.SetValidator(ctx, val)

	// delegator2 delegate 150 CET
	del2Msg := staking.NewMsgDelegate(del2Acc.Address, valAddr, dex.NewCetCoin(150))
	del2Tx := testutil.NewStdTxBuilder("c1").
		Msgs(del2Msg).GasAndFee(1000000, 100).AccNumSeqKey(2, 0, del2Key).Build()
	del2Result := app.Deliver(del2Tx)
	require.Equal(t, sdk.CodeOK, del2Result.Code)

	// assertions
	val, _ = app.stakingKeeper.GetValidator(ctx, valAddr)
	del0, _ := app.stakingKeeper.GetDelegation(ctx, valAcc.Address, valAddr)
	del1, _ := app.stakingKeeper.GetDelegation(ctx, del1Acc.Address, valAddr)
	del2, _ := app.stakingKeeper.GetDelegation(ctx, del2Acc.Address, valAddr)
	require.Equal(t, sdk.NewInt(300), val.Tokens)
	require.Equal(t, sdk.NewDec(400), val.DelegatorShares)
	require.Equal(t, sdk.NewDec(100), del0.Shares)
	require.Equal(t, sdk.NewDec(100), del1.Shares)
	require.Equal(t, sdk.NewDec(200), del2.Shares)
}

func TestSlashTokensToCommunityPool(t *testing.T) {
	// prepare accounts
	valKey, valAcc := testutil.NewBaseAccount(1e9, 0, 0)
	valAddr := sdk.ValAddress(valAcc.Address)

	// init app
	app := initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, valAcc)
		genState.StakingXData.Params.MinSelfDelegation = sdk.NewInt(1e8)
	})

	//begin block at height 1
	//note: context need to be updated after beginblock
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})
	ctx := app.NewContext(false, abci.Header{Height: 1})
	app.stakingKeeper.SetPool(ctx, staking.Pool{
		NotBondedTokens: sdk.NewInt(1e9),
		BondedTokens:    sdk.ZeroInt(),
	})

	// create validator & self delegate 1 CET
	createValMsg := testutil.NewMsgCreateValidatorBuilder(valAddr, valAcc.PubKey).
		MinSelfDelegation(1e8).SelfDelegation(1e8).
		Build()
	createValTx := testutil.NewStdTxBuilder("c1").
		Msgs(createValMsg).GasAndFee(10000, 1).AccNumSeqKey(0, 0, valKey).Build()
	createValResult := app.Deliver(createValTx)
	require.Equal(t, sdk.CodeOK, createValResult.Code)
	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	//create double sign evidence for validator at height 1
	evidences := []abci.Evidence{
		{
			Type:             types.ABCIEvidenceTypeDuplicateVote,
			Validator:        abci.Validator{Address: valAddr, Power: 100},
			Height:           1,
			TotalVotingPower: 100,
		},
	}

	//begin block at height 2 with evidences
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}, ByzantineValidators: evidences})
	ctx = app.NewContext(false, abci.Header{Height: 2})
	app.EndBlock(abci.RequestEndBlock{Height: 2})
	app.Commit()

	//validator should be slashed
	validator, _ := app.stakingKeeper.GetValidator(ctx, valAddr)
	require.Equal(t, sdk.NewInt(95e6), validator.GetTokens())

	//slash tokens should be put into communityPool
	communityPool := app.distrKeeper.GetFeePool(ctx).CommunityPool
	require.Equal(t, sdk.NewDecCoins(dex.NewCetCoins(5e6+1)), communityPool)
}
