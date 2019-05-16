package app

import (
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

	"github.com/coinexchain/dex/denoms"
	"github.com/coinexchain/dex/testutil"
)

func TestMemo(t *testing.T) {
	// genesis state
	toAddr := sdk.AccAddress([]byte("from"))
	key, _, fromAddr := testutil.KeyPubAddr()
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: denoms.NewCetCoins(1000)}
	genAcc := gaia_app.NewGenesisAccount(&acc0)

	genState := NewDefaultGenesisState()
	genState.Accounts = append(genState.Accounts, genAcc)

	// app
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	app := NewCetChainApp(logger, db, nil, true, 10000)

	// init chain
	genStateBytes, _ := app.cdc.MarshalJSON(genState)
	app.InitChain(abci.RequestInitChain{ChainId: "c1", AppStateBytes: genStateBytes})

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins := denoms.NewCetCoins(100)
	msg := bank.NewMsgSend(fromAddr, toAddr, coins)
	fee := auth.NewStdFee(1000000, denoms.NewCetCoins(100))
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).Fee(fee).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
}
