package bankx

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx     sdk.Context
	ak      auth.AccountKeeper
	pk      params.Keeper
	bk      bank.Keeper
	bxk     Keeper
	axk     authx.AccountXKeeper
	ask     ExpectedAssetStatusKeeper
	handler sdk.Handler
}

func (input testInput) handle(msg sdk.Msg) sdk.Result {
	return input.handler(input.ctx, msg)
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	authxKey := sdk.NewKVStoreKey(authx.StoreKey)
	fckKey := sdk.NewKVStoreKey(auth.FeeStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, fckKey)
	axk := authx.NewKeeper(cdc, authxKey, paramsKeeper.Subspace(authx.DefaultParamspace))
	bxkKeeper := NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck, fakeAssetStatusKeeper{})

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, DefaultParams())

	handler := NewHandler(bxkKeeper)
	return testInput{ctx: ctx, ak: ak, pk: paramsKeeper, bk: bk, bxk: bxkKeeper, axk: axk, ask: fakeAssetStatusKeeper{}, handler: handler}
}

func TestHandlerMsgSend(t *testing.T) {
	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")

	fromAccount := input.ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(100000000)
	_ = fromAccount.SetCoins(oneCoins)

	input.ak.SetAccount(input.ctx, fromAccount)
	input.axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	input.handle(msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found := input.axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, false, found)
	require.Equal(t, sdk.NewInt(100000000), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))

	fee := input.bxk.GetParam(input.ctx).LockCoinsFee
	_ = fromAccount.SetCoins(dex.NewCetCoins(1000000000 + fee*2))
	input.ak.SetAccount(input.ctx, fromAccount)

	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee*2), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))

	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(800000000+fee*2), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))

	newMsg := MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 1}
	input.handle(newMsg)
	aux, _ := input.axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(700000000+fee), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)

	newMsg2 := MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(newMsg2)
	aux, _ = input.axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(600000000), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[1].Coin.Amount)
	require.Equal(t, int64(2), aux.LockedCoins[1].UnlockTime)

}

func TestHandlerMsgSendUnlockFirst(t *testing.T) {
	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")
	fromAccount := input.ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)
	fee := input.bxk.GetParam(input.ctx).LockCoinsFee
	Coins := dex.NewCetCoins(1000000000 + fee*2)
	_ = fromAccount.SetCoins(Coins)
	input.ak.SetAccount(input.ctx, fromAccount)
	input.axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found := input.axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, true, found)
	require.Equal(t, sdk.NewInt(100000000+fee), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))

	msgSend2 := MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(msgSend2)
	require.Equal(t, sdk.NewInt(800000000), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found2 := input.axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, true, found2)
	require.Equal(t, sdk.NewInt(100000000+fee*2), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
}

func TestHandleMsgSetMemoRequiredAccountNotExisted(t *testing.T) {
	input := setupTestInput()

	msg := NewMsgSetTransferMemoRequired(testutil.ToAccAddress("xxx"), true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountNotActivated(t *testing.T) {
	input := setupTestInput()

	addr := testutil.ToAccAddress("myaddr")
	//accX := authx.NewAccountXWithAddress(addr)
	//input.axk.SetAccountX(input.ctx, accX)

	msg := NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountOK(t *testing.T) {
	input := setupTestInput()

	addr := testutil.ToAccAddress("myaddr")
	acc := auth.NewBaseAccountWithAddress(addr)
	input.ak.SetAccount(input.ctx, &acc)

	//accX, _ = input.axk.GetAccountX(input.ctx, addr)
	//require.Equal(t, false, accX.MemoRequired)

	msg := NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodeOK, result.Code)

	accX, _ := input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t, true, accX.MemoRequired)
}
