package app

import (
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/store/errors"
	"github.com/stretchr/testify/require"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	gaia_app "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func newApp() *CetChainApp {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	return NewCetChainApp(logger, db, nil, true, 10000)
}

func initApp(accs ...auth.BaseAccount) *CetChainApp {
	app := newApp()

	// genesis state
	genState := NewDefaultGenesisState()
	for _, acc := range accs {
		genAcc := gaia_app.NewGenesisAccount(&acc)
		genState.Accounts = append(genState.Accounts, genAcc)
	}

	// init chain
	genStateBytes, _ := app.cdc.MarshalJSON(genState)
	app.InitChain(abci.RequestInitChain{ChainId: "c1", AppStateBytes: genStateBytes})

	return app
}

func TestSend(t *testing.T) {
	toAddr := sdk.AccAddress([]byte("addr"))
	key, _, fromAddr := testutil.KeyPubAddr()
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initApp(acc0)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins := dex.NewCetCoins(100)
	msg := bank.NewMsgSend(fromAddr, toAddr, coins)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
}

func TestMemo(t *testing.T) {
	key, _, addr := testutil.KeyPubAddr()
	acc0 := auth.BaseAccount{Address: addr, Coins: dex.NewCetCoins(1000)}

	// app
	app := initApp(acc0)

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
	msgSend := bank.NewMsgSend(addr, addr, coins)
	tx2 := testutil.NewStdTxBuilder("c1").
		Msgs(msgSend).Fee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result2 := app.Deliver(tx2)
	require.Equal(t, dex.CodeMemoMissing, result2.Code)
}
