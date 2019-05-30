package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

type genesisStateCallback func(state GenesisState)

func newApp() *CetChainApp {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	return NewCetChainApp(logger, db, nil, true, 10000)
}

func initApp(acc auth.BaseAccount, cb genesisStateCallback) *CetChainApp {
	app := newApp()

	// genesis state
	genState := NewDefaultGenesisState()
	genAcc := NewGenesisAccount(&acc)
	genState.Accounts = append(genState.Accounts, genAcc)
	if cb != nil {
		cb(genState)
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
	app := initApp(acc0, nil)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, time.Now().Unix()+10000)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, bankx.CodeCetCantBeLocked, result.Code)

	msg = bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx = testutil.NewStdTxBuilder("c1").
		Msgs(msg).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result = app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
}

func TestMemo(t *testing.T) {
	key, _, addr := testutil.KeyPubAddr()
	acc0 := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initApp(acc0, nil)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	msgSetMemoRequired := bankx.NewMsgSetTransferMemoRequired(addr, true)
	tx1 := testutil.NewStdTxBuilder("c1").
		Msgs(msgSetMemoRequired).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()
	result1 := app.Deliver(tx1)
	require.Equal(t, errors.CodeOK, result1.Code)

	coins := dex.NewCetCoins(100)
	msgSend := bankx.NewMsgSend(addr, addr, coins, 0)
	tx2 := testutil.NewStdTxBuilder("c1").
		Msgs(msgSend).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result2 := app.Deliver(tx2)
	require.Equal(t, bankx.CodeMemoMissing, result2.Code)
}

func TestGasFeeDeductedWhenTxFailed(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 10000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initApp(acc0, nil)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(100000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeInsufficientCoins, result.Code)

	// end block & commit
	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	// check coins
	ctx := app.NewContext(true, abci.Header{})
	require.Equal(t, int64(10000000000-100),
		app.accountKeeper.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet").Int64())
}

func TestSendFromIncentiveAddr(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	fromAddr := incentive.IncentiveCoinsAccAddr
	coins := dex.NewCetCoinsE8(100)
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}

	// app
	app := initApp(acc0, nil)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins = dex.NewCetCoins(1000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	msgs := make([]sdk.Msg, 1)
	msgs[0] = msg
	tx := auth.StdTx{
		Msgs: msgs,
	}

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeUnauthorized, result.Code)
}

func TestMinSelfDelegation(t *testing.T) {

}