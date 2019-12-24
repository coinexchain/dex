package app

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/cet-sdk/modules/bankx"
	"github.com/coinexchain/cet-sdk/testutil"
	dex "github.com/coinexchain/cet-sdk/types"
)

func TestAccount2UnconfirmedTx(t *testing.T) {
	_, _, toAddr := testutil.KeyPubAddr()
	key, _, fromAddr := testutil.KeyPubAddr()
	key2, _, fromAddr2 := testutil.KeyPubAddr()
	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 30000000000), sdk.NewInt64Coin("eth", 100000000000))
	acc0 := auth.BaseAccount{Address: fromAddr, Coins: coins}
	acc1 := auth.BaseAccount{Address: fromAddr2, Coins: coins}
	// app
	app := initAppWithBaseAccounts(acc0, acc1)
	app.enableUnconfirmedLimit = true
	app.account2UnconfirmedTx.limitTime = 100
	// begin block
	now := time.Now()
	header := abci.Header{Height: 1, Time: now}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	// build tx
	coins = dex.NewCetCoins(1000000000)
	msg := bankx.NewMsgSend(fromAddr, toAddr, coins, now.Unix()+10000)
	tx := newStdTxBuilder().
		Msgs(msg).GasAndFee(600000, 1200000000).AccNumSeqKey(0, 0, key).Build()

	//simple check tx
	txBytes, _ := auth.DefaultTxEncoder(app.cdc)(tx)
	hashID := tmtypes.Tx(txBytes).Hash()
	exist := app.account2UnconfirmedTx.Lookup(fromAddr, hashID, header.Time.Unix())
	require.Equal(t, exist, NoTxExist)
	app.account2UnconfirmedTx.Add(fromAddr, hashID, header.Time.Unix())

	//deliver tx
	result := app.Deliver(tx)
	require.Equal(t, errors.CodeOK, result.Code)
	acc := app.account2UnconfirmedTx.removeList[0]
	require.True(t, bytes.Equal(acc, fromAddr))

	//build another address tx
	tx2 := newStdTxBuilder().
		Msgs(msg).GasAndFee(600000, 1200000000).AccNumSeqKey(0, 0, key2).Build()
	txBytes2, _ := auth.DefaultTxEncoder(app.cdc)(tx2)
	hashIDAnother := tmtypes.Tx(txBytes2).Hash()
	exist = app.account2UnconfirmedTx.Lookup(fromAddr2, hashIDAnother, header.Time.Unix())
	require.Equal(t, exist, NoTxExist)
	app.account2UnconfirmedTx.Add(fromAddr, hashIDAnother, header.Time.Unix())

	//build another same address tx
	msg = bankx.NewMsgSend(fromAddr, toAddr, coins, 0)
	tx = newStdTxBuilder().
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 1, key).Build()

	//check should failed for unconfirmed tx already in map
	txBytes, _ = auth.DefaultTxEncoder(app.cdc)(tx)
	hashID2 := tmtypes.Tx(txBytes).Hash()
	require.NotEqual(t, hashID, hashID2)
	exist = app.account2UnconfirmedTx.Lookup(fromAddr, hashID2, header.Time.Unix())
	require.Equal(t, exist, OtherTxExist)

	//end block
	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()
	require.Equal(t, len(app.account2UnconfirmedTx.auMap), 0)

	//next block
	header = abci.Header{Height: 2}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	require.Equal(t, len(app.account2UnconfirmedTx.removeList), 0)

	tx = newStdTxBuilder().
		Msgs(msg).GasAndFee(600000, 1200000000).AccNumSeqKey(0, 1, key).Build()

	//simple check tx
	txBytes, _ = auth.DefaultTxEncoder(app.cdc)(tx)
	hashID3 := tmtypes.Tx(txBytes).Hash()
	exist = app.account2UnconfirmedTx.Lookup(fromAddr, hashID3, header.Time.Unix())
	require.Equal(t, exist, NoTxExist)
	app.account2UnconfirmedTx.Add(fromAddr, hashID3, header.Time.Unix())
}
