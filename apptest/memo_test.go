package apptest

import (
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/denoms"
)

func TestMemo(t *testing.T) {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	app := app.NewCetChainApp(logger, db, nil, true, 10000)

	app.InitChain(abci.RequestInitChain{AppStateBytes:[]byte("{}")})

	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	fromAddr := sdk.AccAddress([]byte("from"))
	toAddr := sdk.AccAddress([]byte("from"))
	coins := denoms.NewCetCoins(100)
	msg := bank.NewMsgSend(fromAddr, toAddr, coins)
	msgs := []sdk.Msg{msg}
	fee := auth.NewStdFee(10000, denoms.NewCetCoins(100))
	sigs := []auth.StdSignature{}
	tx := auth.NewStdTx(msgs, fee, sigs, "")
	app.Deliver(tx)
}
