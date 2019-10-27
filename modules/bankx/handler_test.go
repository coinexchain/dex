package bankx_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	bx "github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

var (
	myaddr        = testutil.ToAccAddress("myaddr")
	fromAddr      = testutil.ToAccAddress("fromaddr")
	toAddr        = testutil.ToAccAddress("toaddr")
	supervisor    = testutil.ToAccAddress("supervisoraddr")
	feeAddr       = sdk.AccAddress(crypto.AddressHash([]byte(auth.FeeCollectorName)))
	owner         = testutil.ToAccAddress("owner")
	forbiddenAddr = testutil.ToAccAddress("forbidden")
)

func defaultContext() (*keeper.Keeper, sdk.Handler, sdk.Context) {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{Time: time.Now()}, false, log.NewNopLogger())
	app.BankxKeeper.SetParams(ctx, bx.DefaultParams())
	app.BankxKeeper.SetSendEnabled(ctx, true)
	handler := bankx.NewHandler(app.BankxKeeper)
	cet, _ := asset.NewToken("cet", "cet", sdk.NewInt(200000000000000), owner,
		false, false, true, true,
		"", "", asset.TestIdentityString)
	_ = app.AssetKeeper.SetToken(ctx, cet)
	_ = app.AssetKeeper.ForbidAddress(ctx, "cet", owner, []sdk.AccAddress{forbiddenAddr})
	return &app.BankxKeeper, handler, ctx
}

func TestHandlerMsgSend(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100000000))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)
	fee := bkx.GetParams(ctx).LockCoinsFeePerDay
	err = bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(1000000000+fee*2))
	require.NoError(t, err)

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee*2), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(800000000+fee*2), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	newMsg := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime + 1}
	handle(ctx, newMsg)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(700000000+fee), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 1, len(bkx.GetLockedCoins(ctx, toAddr)))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, lockFreeTime+1, bkx.GetLockedCoins(ctx, toAddr)[0].UnlockTime)

	newMsg2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime + 2}
	handle(ctx, newMsg2)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(600000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 2, len(bkx.GetLockedCoins(ctx, toAddr)))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[1].Coin.Amount)
	require.Equal(t, lockFreeTime+2, bkx.GetLockedCoins(ctx, toAddr)[1].UnlockTime)

	newMsg3 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime}
	handle(ctx, newMsg3)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(500000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 3, len(bkx.GetLockedCoins(ctx, toAddr)))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[2].Coin.Amount)
	require.Equal(t, lockFreeTime, bkx.GetLockedCoins(ctx, toAddr)[2].UnlockTime)
}

func TestHandlerMsgSendFail(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100000000))
	require.NoError(t, err)

	bkx.SetSendEnabled(ctx, false)
	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	res := handle(ctx, msgSend)
	require.Equal(t, bank.CodeSendDisabled, res.Code)

	bkx.SetSendEnabled(ctx, true)
	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(200000000), UnlockTime: 0}
	res = handle(ctx, msgSend)
	require.Equal(t, sdk.CodeInsufficientCoins, res.Code)

}

func TestHandlerMsgSendUnlockFirst(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	fee := bkx.GetParams(ctx).LockCoinsFeePerDay
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(1000000000+fee*2))
	require.NoError(t, err)

	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime + 2}
	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	msgSend2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime + 2}
	handle(ctx, msgSend2)
	require.Equal(t, sdk.NewInt(800000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
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

	err := bkx.AddCoins(ctx, myaddr, sdk.Coins{})
	require.NoError(t, err)

	msg := bx.NewMsgSetTransferMemoRequired(myaddr, true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodeOK, result.Code)
	require.Equal(t, true, bkx.GetMemoRequired(ctx, myaddr))
}

func TestUnlockQueueNotAppend(t *testing.T) {
	bkx, handle, ctx := defaultContext()

	fee := bkx.GetParams(ctx).LockCoinsFeePerDay
	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)

	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100000000+fee))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: lockFreeTime + 1000}
	handle(ctx, msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, 0, len(bkx.GetLockedCoins(ctx, toAddr)))
}

func TestLockedCoinFee(t *testing.T) {
	bkx, handle, ctx := defaultContext()

	fee := bkx.GetParams(ctx).LockCoinsFeePerDay
	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)

	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(5e8+fee*3))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(1e8)}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(4e8+fee*3), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(1e8), UnlockTime: lockFreeTime}
	handle(ctx, msgSend)
	// no fee
	require.Equal(t, sdk.NewInt(3e8+fee*3), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 1, len(bkx.GetLockedCoins(ctx, toAddr)))

	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(1e8), UnlockTime: lockFreeTime + 24*3600}
	handle(ctx, msgSend)
	// 1 day fee
	require.Equal(t, sdk.NewInt(2e8+fee*2), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8+fee), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 2, len(bkx.GetLockedCoins(ctx, toAddr)))

	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(1e8), UnlockTime: lockFreeTime + 24*3600 + 1}
	handle(ctx, msgSend)
	// 2 day fees
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8+fee*3), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, 3, len(bkx.GetLockedCoins(ctx, toAddr)))
}

func TestHandlerMsgMultiSend(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	_ = bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(1000000000))
	_ = bkx.AddCoins(ctx, myaddr, dex.NewCetCoins(1000000000))

	coins := sdk.NewCoins(sdk.NewInt64Coin("cet", 300000000))
	in := []bank.Input{bank.NewInput(fromAddr, coins), bank.NewInput(myaddr, coins)}
	out := []bank.Output{bank.NewOutput(toAddr, coins), bank.NewOutput(toAddr, coins)}
	msg := bankx.NewMsgMultiSend(in, out)

	bkx.SetSendEnabled(ctx, false)
	handle(ctx, msg)
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))

	bkx.SetSendEnabled(ctx, true)
	handle(ctx, msg)
	require.Equal(t, sdk.NewInt(700000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(700000000), bkx.GetCoins(ctx, myaddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(500000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	in = []bank.Input{bank.NewInput(fromAddr, coins), bank.NewInput(forbiddenAddr, coins)}
	msg = bankx.NewMsgMultiSend(in, out)
	handle(ctx, msg)
	require.Equal(t, sdk.NewInt(500000000).String(), bkx.GetCoins(ctx, toAddr).AmountOf("cet").String())

	newAddr := testutil.ToAccAddress("newAddr")
	invalid := sdk.NewCoins(sdk.NewInt64Coin("cet", 1000))
	in = []bank.Input{bank.NewInput(fromAddr, invalid), bank.NewInput(myaddr, invalid)}
	out = []bank.Output{bank.NewOutput(toAddr, invalid), bank.NewOutput(newAddr, invalid)}
	msg = bankx.NewMsgMultiSend(in, out)
	res := handle(ctx, msg)
	require.False(t, res.IsOK())
}

func TestHandleMsgSupervisedSend(t *testing.T) {
	bkx, handle, ctx := defaultContext()

	//fee := bkx.GetParams(ctx).LockCoinsFeePerDay
	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)

	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100*1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, toAddr, dex.NewCetCoins(1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, supervisor, dex.NewCetCoins(1e8))
	require.NoError(t, err)

	// create -> return
	msgSend := bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.Create}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(97e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.Return}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(99e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(2e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	// create -> unlock by sender
	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.Create}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(96e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(1e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(2e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.EarlierUnlockBySender}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(96e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(3e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(3e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	// create -> unlock by supervisor
	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.Create}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(93e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(3e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(3e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: supervisor, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.EarlierUnlockBySupervisor}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(93e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(5e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(4e8), bkx.GetCoins(ctx, supervisor).AmountOf("cet"))

	// no supervisor
	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: nil, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.Create}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(90e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(5e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))

	msgSend = bankx.MsgSupervisedSend{FromAddress: fromAddr, Supervisor: nil, ToAddress: toAddr,
		Amount: dex.NewCetCoin(3e8), UnlockTime: lockFreeTime, Reward: 1e8, Operation: bankx.EarlierUnlockBySender}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(90e8), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(8e8), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
}

func TestHandleMsgSupervisedSendException(t *testing.T) {
	bkx, handle, ctx := defaultContext()

	params := bkx.GetParams(ctx)
	params.LockCoinsFeePerDay = 1e15
	bkx.SetParams(ctx, params)
	lockFreeTime := ctx.BlockHeader().Time.Unix() + bkx.GetParams(ctx).LockCoinsFreeTime/int64(time.Second)
	govModuleAcc := supply.NewEmptyModuleAccount("gov")

	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100*1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, toAddr, dex.NewCetCoins(1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, supervisor, dex.NewCetCoins(1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, forbiddenAddr, dex.NewCetCoins(100*1e8))
	require.NoError(t, err)
	err = bkx.AddCoins(ctx, govModuleAcc.Address, dex.NewCetCoins(100*1e8))
	require.NoError(t, err)

	var msgSend bankx.MsgSupervisedSend
	var ret sdk.Result
	var emptyAddr sdk.AccAddress
	newAddr := testutil.ToAccAddress("newaddr")
	amt := dex.NewCetCoin(3e8)

	bkx.SetSendEnabled(ctx, false)
	msgSend = bx.NewMsgSupervisedSend(fromAddr, supervisor, toAddr, amt, lockFreeTime, 1e8, bankx.Create)
	ret = handle(ctx, msgSend)
	require.False(t, ret.IsOK())
	bkx.SetSendEnabled(ctx, true)

	now := time.Now()
	header := abci.Header{Time: now, Height: 10}
	ctx = ctx.WithBlockHeader(header)

	fmt.Printf("max time < :%v\n", math.MaxInt64/int64(time.Second))
	fmt.Printf("count:%v\n", math.MaxInt64/int64(24*time.Hour))

	cases := []struct {
		isOK bool
		code sdk.CodeType
		msg  bankx.MsgSupervisedSend
	}{
		{false, sdk.CodeUnknownAddress,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, emptyAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, sdk.CodeUnknownAddress,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, newAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, sdk.CodeUnauthorized,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, govModuleAcc.Address, amt, lockFreeTime, 1e8, bankx.Create)},
		{true, sdk.CodeOK,
			bx.NewMsgSupervisedSend(fromAddr, emptyAddr, toAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, sdk.CodeUnknownAddress,
			bx.NewMsgSupervisedSend(fromAddr, newAddr, toAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, sdk.CodeUnauthorized,
			bx.NewMsgSupervisedSend(fromAddr, govModuleAcc.Address, toAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, sdk.CodeInsufficientCoins,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, toAddr, dex.NewCetCoin(1000e8), lockFreeTime, 1e8, bankx.Create)},
		{false, bx.CodeInvalidUnlockTime,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, toAddr, amt, ctx.BlockHeader().Time.Unix()-1, 1e8, bankx.Create)},
		{false, bx.CodeInvalidUnlockTime,
			bx.NewMsgSupervisedSend(fromAddr, supervisor, toAddr, amt, ctx.BlockHeader().Time.Unix()+math.MaxInt64/int64(1e9), 1e8, bankx.Create)},
		{false, bx.CodeTokenForbiddenByOwner,
			bx.NewMsgSupervisedSend(forbiddenAddr, supervisor, toAddr, amt, lockFreeTime, 1e8, bankx.Create)},
		{false, bx.CodeLockedCoinNotFound,
			bx.NewMsgSupervisedSend(forbiddenAddr, supervisor, toAddr, amt, lockFreeTime, 1e8, bankx.Return)},
	}

	for _, tc := range cases {
		ret = handle(ctx, tc.msg)
		require.Equal(t, tc.isOK, ret.IsOK())
		require.Equal(t, tc.code, ret.Code)
	}
}
