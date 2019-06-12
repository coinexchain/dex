package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func TestGasFeeDeductedWhenTxFailed(t *testing.T) {
	// acc & app
	key, acc := testutil.NewBaseAccount(10000000000, 0, 0)
	app := initAppWithBaseAccounts(acc)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins := dex.NewCetCoins(100000000000)
	toAddr := sdk.AccAddress([]byte("addr"))
	msg := bankx.NewMsgSend(acc.Address, toAddr, coins, 0)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(1000000, 100).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, sdk.CodeInsufficientCoins, result.Code)

	// end block & commit
	app.EndBlock(abci.RequestEndBlock{Height: 1})
	app.Commit()

	// check coins
	ctx := app.NewContext(true, abci.Header{})
	require.Equal(t, int64(10000000000-100),
		app.accountKeeper.GetAccount(ctx, acc.Address).GetCoins().AmountOf("cet").Int64())
}

func TestMinGasPriceLimit(t *testing.T) {
	// acc & app
	key, acc := testutil.NewBaseAccount(1e10, 0, 0)
	app := initAppWithBaseAccounts(acc)

	// begin block
	header := abci.Header{Height: 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// deliver tx
	coins := dex.NewCetCoins(1e8)
	toAddr := sdk.AccAddress([]byte("addr"))
	msg := bankx.NewMsgSend(acc.Address, toAddr, coins, 0)
	tx := testutil.NewStdTxBuilder("c1").
		Msgs(msg).GasAndFee(10000000000, 1).AccNumSeqKey(0, 0, key).Build()

	result := app.Deliver(tx)
	require.Equal(t, authx.CodeGasPriceTooLow, result.Code)
}
