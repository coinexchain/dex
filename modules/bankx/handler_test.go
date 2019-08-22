package bankx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	cetapp "github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	bx "github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

var (
	myaddr   = testutil.ToAccAddress("myaddr")
	fromAddr = testutil.ToAccAddress("fromaddr")
	toAddr   = testutil.ToAccAddress("toaddr")
	feeAddr  = sdk.AccAddress(crypto.AddressHash([]byte(auth.FeeCollectorName)))
	owner    = testutil.ToAccAddress("owner")
)

func defaultContext() (*keeper.Keeper, sdk.Handler, sdk.Context) {
	app := cetapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{}, false, log.NewNopLogger())
	app.BankxKeeper.SetParams(ctx, bx.DefaultParams())
	app.BankxKeeper.Bk.SetSendEnabled(ctx, true)
	handler := bankx.NewHandler(app.BankxKeeper)
	cet, _ := asset.NewToken("cet", "cet", sdk.NewInt(200000000000000), owner,
		false, false, false, false,
		"", "", "")
	_ = app.AssetKeeper.SetToken(ctx, cet)
	return &app.BankxKeeper, handler, ctx
}

func TestHandlerMsgSend(t *testing.T) {

	bkx, handle, ctx := defaultContext()

	fromAccount := bkx.Ak.NewAccountWithAddress(ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(100000000)
	_ = fromAccount.SetCoins(oneCoins)

	bkx.Ak.SetAccount(ctx, fromAccount)
	bkx.Axk.SetAccountX(ctx, fromAccountX)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(0).String(), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet").String())
	require.Equal(t, sdk.NewInt(0).String(), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet").String())
	_, found := bkx.Axk.GetAccountX(ctx, toAddr)
	require.Equal(t, false, found)
	require.Equal(t, sdk.NewInt(100000000), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))

	fee := bkx.GetParams(ctx).LockCoinsFee
	_ = fromAccount.SetCoins(dex.NewCetCoins(1000000000 + fee*2))
	bkx.Ak.SetAccount(ctx, fromAccount)

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee*2), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(800000000+fee*2), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))

	newMsg := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 1}
	handle(ctx, newMsg)
	aux, _ := bkx.Axk.GetAccountX(ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(700000000+fee), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)

	newMsg2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, newMsg2)
	aux, _ = bkx.Axk.GetAccountX(ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(600000000), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[1].Coin.Amount)
	require.Equal(t, int64(2), aux.LockedCoins[1].UnlockTime)
}

func TestHandlerMsgSendFail(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	fromAccount := bkx.Ak.NewAccountWithAddress(ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(100000000)
	_ = fromAccount.SetCoins(oneCoins)

	bkx.Ak.SetAccount(ctx, fromAccount)
	bkx.Axk.SetAccountX(ctx, fromAccountX)

	bkx.Bk.SetSendEnabled(ctx, false)
	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	res := handle(ctx, msgSend)
	require.Equal(t, bank.CodeSendDisabled, res.Code)

	bkx.Bk.SetSendEnabled(ctx, true)
	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(200000000), UnlockTime: 0}
	res = handle(ctx, msgSend)
	require.Equal(t, sdk.CodeInsufficientCoins, res.Code)

}

func TestHandlerMsgSendUnlockFirst(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	fromAccount := bkx.Ak.NewAccountWithAddress(ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)
	fee := bkx.GetParams(ctx).LockCoinsFee
	Coins := dex.NewCetCoins(1000000000 + fee*2)
	_ = fromAccount.SetCoins(Coins)
	bkx.Ak.SetAccount(ctx, fromAccount)
	bkx.Axk.SetAccountX(ctx, fromAccountX)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found := bkx.Axk.GetAccountX(ctx, toAddr)
	require.Equal(t, true, found)
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))

	msgSend2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, msgSend2)
	require.Equal(t, sdk.NewInt(800000000), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found2 := bkx.Axk.GetAccountX(ctx, toAddr)
	require.Equal(t, true, found2)
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.Ak.GetAccount(ctx, feeAddr).GetCoins().AmountOf("cet"))
}

func TestHandleMsgSetMemoRequiredAccountNotExisted(t *testing.T) {
	_, handle, ctx := defaultContext()

	msg := bankx.NewMsgSetTransferMemoRequired(testutil.ToAccAddress("xxx"), true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountNotActivated(t *testing.T) {
	_, handle, ctx := defaultContext()

	msg := bankx.NewMsgSetTransferMemoRequired(myaddr, true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountOK(t *testing.T) {
	bkx, handle, ctx := defaultContext()

	acc := auth.NewBaseAccountWithAddress(myaddr)
	bkx.Ak.SetAccount(ctx, &acc)

	msg := bx.NewMsgSetTransferMemoRequired(myaddr, true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodeOK, result.Code)

	accX, _ := bkx.Axk.GetAccountX(ctx, myaddr)
	require.Equal(t, true, accX.MemoRequired)
}

func TestUnlockQueueNotAppend(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	fromAccount := bkx.Ak.NewAccountWithAddress(ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(10100000000)
	_ = fromAccount.SetCoins(oneCoins)

	bkx.Ak.SetAccount(ctx, fromAccount)
	bkx.Axk.SetAccountX(ctx, fromAccountX)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 10000}
	handle(ctx, msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), bkx.Ak.GetAccount(ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.Ak.GetAccount(ctx, toAddr).GetCoins().AmountOf("cet"))
}
