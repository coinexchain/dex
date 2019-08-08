package app

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	types2 "github.com/coinexchain/dex/modules/distributionx/types"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

const testChainID = "c1"

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

type genesisStateCallback func(state *GenesisState)

func newStdTxBuilder() *testutil.StdTxBuilder {
	return testutil.NewStdTxBuilder(testChainID)
}

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
	var amount int64
	for _, acc := range accs {
		genAcc := genaccounts.NewGenesisAccount(&acc)
		genState.Accounts = append(genState.Accounts, genAcc)
		amount = amount + acc.Coins.AmountOf(dex.CET).Int64()
	}

	addAccountForDanglingCET(amount, genState)
}
func addModuleAccounts(genState *GenesisState) {
	maccs := []*supply.ModuleAccount{
		supply.NewEmptyModuleAccount(auth.FeeCollectorName),
		supply.NewEmptyModuleAccount(distribution.ModuleName),
		supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking),
		supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking),
		supply.NewEmptyModuleAccount(gov.ModuleName, supply.Burner),
		supply.NewEmptyModuleAccount(authx.ModuleName),
		supply.NewEmptyModuleAccount(asset.ModuleName, supply.Burner, supply.Minter),
	}
	for _, macc := range maccs {
		genMacc, err := genaccounts.NewGenesisAccountI(*macc)
		if err != nil {
			genState.Accounts = append(genState.Accounts, genMacc)
		}
	}
}

func addAccountForDanglingCET(amount int64, genState *GenesisState) {
	accAmount := cetToken().GetTotalSupply().Int64() - amount
	if accAmount > 0 {
		_, acc := testutil.NewBaseAccount(accAmount, 1, 0)
		genAcc := genaccounts.NewGenesisAccount(&acc)
		genState.Accounts = append(genState.Accounts, genAcc)
	}
}

func initApp(cb genesisStateCallback) *CetChainApp {
	app := newApp()

	// genesis state
	genState := NewDefaultGenesisState()

	cetToken := cetToken()
	genState.AssetData.Tokens = append(genState.AssetData.Tokens, cetToken)
	genState.StakingData.Params.BondDenom = dex.DefaultBondDenom

	genState.AuthXData.Params.MinGasPriceLimit = sdk.MustNewDecFromStr("0.00000001")
	if cb != nil {
		cb(&genState)
	}

	// init chain
	genStateBytes, _ := app.cdc.MarshalJSON(genState)
	app.InitChain(abci.RequestInitChain{ChainId: testChainID, AppStateBytes: genStateBytes})

	return app
}

func initAppWithAccounts(accs ...auth.BaseAccount) *CetChainApp {
	return initApp(func(genState *GenesisState) {
		addGenesisAccounts(genState, accs...)
		addModuleAccounts(genState)
	})
}
func cetToken() asset.Token {
	cetOwnerAddr, _ := sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	return &asset.BaseToken{
		Name:             "CoinEx Chain Native Token",
		Symbol:           "cet",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            cetOwnerAddr,
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
	}
}

func TestSend(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000), sdk.NewInt64Coin("eth", 100000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, time.Now().Unix()+10000)
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeType(0), result.Code)

	msg = bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx = newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 1, key).Build()

	result = app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
}
func TestBankSend(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000), sdk.NewInt64Coin("eth", 100000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	msg := bank.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      coins,
	}
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)

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
	tx1 := newStdTxBuilder().
		Msgs(msgSetMemoRequired).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()
	result1 := app.Deliver(tx1)
	require.Equal(t, errors.CodeOK, result1.Code)

	coins := dex.NewCetCoins(100)
	msgSend := bankx.NewMsgSend(addr, addr, coins, 0)
	tx2 := newStdTxBuilder().
		Msgs(msgSend).GasAndFee(1000000, 100).AccNumSeqKey(0, 1, key).Build()

	result2 := app.Deliver(tx2)
	require.Equal(t, bankx.CodeMemoMissing, result2.Code)
}

func TestSendFromIncentiveAddr(t *testing.T) {
	key, _, fromAddr := testutil.KeyPubAddr()
	incentive.PoolAddr = fromAddr
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
	tx := newStdTxBuilder().
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
		genState.Accounts = append(genState.Accounts, genaccounts.NewGenesisAccount(&acc0))
		genState.StakingXData.Params.MinSelfDelegation = sdk.NewInt(500)

		addAccountForDanglingCET(1000, genState)
	})

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver
	msg := testutil.NewMsgCreateValidatorBuilder(val0, pubKey0).
		MinSelfDelegation(400).SelfDelegation(450).
		Build()
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key0).Build()

	result := app.Deliver(tx)
	//require.Nil(t, result.Codespace)
	require.Equal(t, stakingx.CodeMinSelfDelegationBelowRequired, result.Code)
}

func TestDelegatorShares(t *testing.T) {
	// prepare accounts
	amountVal := cetToken().GetTotalSupply().Int64() - 20000
	valKey, valAcc := testutil.NewBaseAccount(amountVal, 0, 0)
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
		Commission("0.05", "0.1", "0.01").
		Build()
	createValTx := newStdTxBuilder().
		Msgs(createValMsg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, valKey).Build()
	createValResult := app.Deliver(createValTx)
	require.Equal(t, sdk.CodeOK, createValResult.Code)

	// delegator1 delegate 100 CET
	del1Msg := staking.NewMsgDelegate(del1Acc.Address, valAddr, dex.NewCetCoin(100))
	del1Tx := newStdTxBuilder().
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
	del2Tx := newStdTxBuilder().
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

	// create validator & self delegate 1 CET
	createValMsg := testutil.NewMsgCreateValidatorBuilder(valAddr, valAcc.PubKey).
		MinSelfDelegation(1e8).SelfDelegation(1e8).
		Commission("0.05", "0.1", "0.01").
		Build()
	createValTx := newStdTxBuilder().
		Msgs(createValMsg).GasAndFee(1000000, 1).AccNumSeqKey(0, 0, valKey).Build()
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
	ctx := app.NewContext(false, abci.Header{Height: 2})
	app.EndBlock(abci.RequestEndBlock{Height: 2})
	app.Commit()

	//validator should be slashed
	validator, _ := app.stakingKeeper.GetValidator(ctx, valAddr)
	require.Equal(t, sdk.NewInt(95e6), validator.GetTokens())

	//slash tokens should be put into communityPool
	communityPool := app.distrKeeper.GetFeePool(ctx).CommunityPool
	require.Equal(t, sdk.NewDecCoins(dex.NewCetCoins(5e6+1)), communityPool)
}

func TestDonateToCommunityPool(t *testing.T) {

	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10e8))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := app.NewContext(false, header)

	//build tx
	coins = dex.NewCetCoins(1e8)
	msg := types2.NewMsgDonateToCommunityPool(fromAddr, coins)
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	// deliver tx
	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeType(0), result.Code)

	//check account
	acc := app.accountKeeper.GetAccount(ctx, fromAddr)
	require.Equal(t, sdk.NewInt(899999900), acc.GetCoins().AmountOf("cet"))

	//check communityPool
	communityPool := app.distrKeeper.GetFeePool(ctx).CommunityPool
	require.True(t, communityPool.AmountOf("cet").Equal(sdk.NewDec(1e8)))
}

func TestTotalSupplyInvariant(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000), sdk.NewInt64Coin("eth", 100000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := app.NewContext(false, header)
	app.crisisKeeper.SetConstantFee(ctx, sdk.NewCoin("cet", sdk.NewInt(1e8)))

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, time.Now().Unix()+10000)
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeType(0), result.Code)

	authxMacc := app.supplyKeeper.GetModuleAccount(ctx, authx.ModuleName)
	require.True(t, authxMacc.GetCoins().Empty())

	msgInv := crisis.NewMsgVerifyInvariant(fromAddr, supply.ModuleName, "total-supply")
	tx = newStdTxBuilder().
		Msgs(msgInv).GasAndFee(1000000, 100).AccNumSeqKey(0, 1, key).Build()

	result = app.Deliver(tx)
	require.Equal(t, sdk.CodeOK, result.Code)

	authxMacc = app.supplyKeeper.GetModuleAccount(ctx, authx.ModuleName)
	require.False(t, authxMacc.GetCoins().Empty())
}

func TestBankMultiSend(t *testing.T) {
	key1, _, fromAddr1 := testutil.KeyPubAddr()
	key2, _, fromAddr2 := testutil.KeyPubAddr()
	_, _, toAddr1 := testutil.KeyPubAddr()
	_, _, toAddr2 := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000))
	acc0 := auth.BaseAccount{Address: fromAddr1, Coins: coins}
	acc1 := auth.BaseAccount{Address: fromAddr2, Coins: coins}

	// app
	app := initAppWithBaseAccounts(acc0, acc1)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000000000)

	msg := bank.MsgMultiSend{
		Inputs: []bank.Input{
			bank.NewInput(fromAddr1, coins),
			bank.NewInput(fromAddr2, coins),
		},
		Outputs: []bank.Output{
			bank.NewOutput(toAddr1, coins),
			bank.NewOutput(toAddr2, coins),
		},
	}

	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).
		AccNumSeqKey(0, 0, key1).AccNumSeqKey(1, 0, key2).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func TestMultiSend(t *testing.T) {
	key1, _, fromAddr1 := testutil.KeyPubAddr()
	key2, _, fromAddr2 := testutil.KeyPubAddr()
	_, _, toAddr1 := testutil.KeyPubAddr()
	_, _, toAddr2 := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000))
	acc0 := auth.BaseAccount{Address: fromAddr1, Coins: coins}
	acc1 := auth.BaseAccount{Address: fromAddr2, Coins: coins}
	acc2 := auth.BaseAccount{Address: toAddr1, Coins: dex.NewCetCoins(0)}

	// app
	app := initAppWithBaseAccounts(acc0, acc1, acc2)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000000000)

	in := []bank.Input{bank.NewInput(fromAddr1, coins), bank.NewInput(fromAddr2, coins)}
	out := []bank.Output{bank.NewOutput(toAddr1, coins), bank.NewOutput(toAddr2, coins)}

	msg := bankx.NewMsgMultiSend(in, out)

	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).
		AccNumSeqKey(0, 0, key1).AccNumSeqKey(1, 0, key2).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeType(0), result.Code)

	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	header = abci.Header{Height: 2}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := app.NewContext(false, header)

	expectedCoins := dex.NewCetCoins(9e8)
	toAcc1 := app.accountKeeper.GetAccount(ctx, toAddr1)
	toAcc2 := app.accountKeeper.GetAccount(ctx, toAddr2)
	require.Equal(t, coins, toAcc1.GetCoins())
	require.Equal(t, expectedCoins, toAcc2.GetCoins())
}

func TestMultiSendMemoRequired(t *testing.T) {
	key1, _, fromAddr1 := testutil.KeyPubAddr()
	key2, _, fromAddr2 := testutil.KeyPubAddr()
	key3, _, toAddr1 := testutil.KeyPubAddr()
	_, _, toAddr2 := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000))
	acc0 := auth.BaseAccount{Address: fromAddr1, Coins: coins}
	acc1 := auth.BaseAccount{Address: fromAddr2, Coins: coins}
	acc2 := auth.BaseAccount{Address: toAddr1, Coins: dex.NewCetCoins(1e8)}

	// app
	app := initAppWithBaseAccounts(acc0, acc1, acc2)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	//deliver tx
	msgSetMemoRequired := bankx.NewMsgSetTransferMemoRequired(toAddr1, true)
	tx1 := newStdTxBuilder().
		Msgs(msgSetMemoRequired).GasAndFee(1000000, 100).AccNumSeqKey(2, 0, key3).Build()
	result1 := app.Deliver(tx1)
	require.Equal(t, errors.CodeOK, result1.Code)

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	in := []bank.Input{bank.NewInput(fromAddr1, coins), bank.NewInput(fromAddr2, coins)}
	out := []bank.Output{bank.NewOutput(toAddr1, coins), bank.NewOutput(toAddr2, coins)}
	msg := bankx.NewMsgMultiSend(in, out)
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).
		AccNumSeqKey(0, 0, key1).AccNumSeqKey(1, 0, key2).Build()
	result := app.Deliver(tx)
	require.Equal(t, bankx.CodeMemoMissing, result.Code)

}
