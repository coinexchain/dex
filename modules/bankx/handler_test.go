package bankx

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/authx"
	types2 "github.com/coinexchain/dex/modules/authx/types"
	bx "github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

type fakeAssetStatusKeeper struct{}

func (k fakeAssetStatusKeeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	return false
}
func (k fakeAssetStatusKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	return false
}
func (k fakeAssetStatusKeeper) UpdateTokenSendLock(ctx sdk.Context, symbol string, amount sdk.Int, lock bool) sdk.Error {
	return nil
}

var myaddr = testutil.ToAccAddress("myaddr")
var feeAddr = sdk.AccAddress(crypto.AddressHash([]byte(auth.FeeCollectorName)))

func defaultContext() (sdk.Context, *codec.Codec, Keeper) {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	bank.RegisterCodec(cdc)
	bx.RegisterCodec(cdc)

	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")

	keySupply := sdk.NewKVStoreKey("supply")
	keyAuth := sdk.NewKVStoreKey("auth")
	keyAuthX := sdk.NewKVStoreKey("authx")
	keyBank := sdk.NewKVStoreKey("bank")
	keyBankx := sdk.NewKVStoreKey("bankx")

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyAuthX, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyBankx, sdk.StoreTypeIAVL, db)
	_ = cms.LoadLatestVersion()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		authx.ModuleName:          nil,
		bank.ModuleName:           nil,
		"bankx":                   nil,
	}

	ask := fakeAssetStatusKeeper{}

	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAuth, paramsKeeper.Subspace("auth"), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace("bank"), "bank", map[string]bool{})
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	axk := authx.NewKeeper(cdc, keyAuthX, paramsKeeper.Subspace("authx"), sk, ak, "")
	bxK := NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, ask, sk, msgqueue.NewProducer())

	return ctx, cdc, bxK
}

type testInput struct {
	ctx     sdk.Context
	bxk     Keeper
	handler sdk.Handler
}

func (input testInput) handle(msg sdk.Msg) sdk.Result {
	return input.handler(input.ctx, msg)
}

func setupTestInput() testInput {

	ctx, _, bxk := defaultContext()
	bxk.Bk.SetSendEnabled(ctx, true)
	bxk.SetParam(ctx, bx.DefaultParams())

	handler := NewHandler(bxk)
	return testInput{ctx: ctx, bxk: bxk, handler: handler}
}

func TestHandlerMsgSend(t *testing.T) {

	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")

	fromAccount := input.bxk.Ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := types2.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(100000000)
	_ = fromAccount.SetCoins(oneCoins)

	input.bxk.Ak.SetAccount(input.ctx, fromAccount)
	input.bxk.Axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	input.handle(msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found := input.bxk.Axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, false, found)
	require.Equal(t, sdk.NewInt(100000000), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))

	fee := input.bxk.GetParam(input.ctx).LockCoinsFee
	_ = fromAccount.SetCoins(dex.NewCetCoins(1000000000 + fee*2))
	input.bxk.Ak.SetAccount(input.ctx, fromAccount)

	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee*2), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))

	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(800000000+fee*2), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))

	newMsg := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 1}
	input.handle(newMsg)
	aux, _ := input.bxk.Axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(700000000+fee), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)

	newMsg2 := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(newMsg2)
	aux, _ = input.bxk.Axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(600000000), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, int64(1), aux.LockedCoins[0].UnlockTime)
	require.Equal(t, sdk.NewInt(100000000), aux.LockedCoins[1].Coin.Amount)
	require.Equal(t, int64(2), aux.LockedCoins[1].UnlockTime)
}

func TestHandlerMsgSendFail(t *testing.T) {
	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")

	fromAccount := input.bxk.Ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := types2.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(100000000)
	_ = fromAccount.SetCoins(oneCoins)

	input.bxk.Ak.SetAccount(input.ctx, fromAccount)
	input.bxk.Axk.SetAccountX(input.ctx, fromAccountX)

	input.bxk.Bk.SetSendEnabled(input.ctx, false)
	msgSend := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	res := input.handle(msgSend)
	require.Equal(t, bank.CodeSendDisabled, res.Code)

	input.bxk.Bk.SetSendEnabled(input.ctx, true)
	msgSend = bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(200000000), UnlockTime: 0}
	res = input.handle(msgSend)
	require.Equal(t, sdk.CodeInsufficientCoins, res.Code)

}

func TestHandlerMsgSendUnlockFirst(t *testing.T) {
	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")
	fromAccount := input.bxk.Ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := types2.NewAccountXWithAddress(fromAddr)
	fee := input.bxk.GetParam(input.ctx).LockCoinsFee
	Coins := dex.NewCetCoins(1000000000 + fee*2)
	_ = fromAccount.SetCoins(Coins)
	input.bxk.Ak.SetAccount(input.ctx, fromAccount)
	input.bxk.Axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found := input.bxk.Axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, true, found)
	require.Equal(t, sdk.NewInt(100000000+fee), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))

	msgSend2 := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	input.handle(msgSend2)
	require.Equal(t, sdk.NewInt(800000000), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
	_, found2 := input.bxk.Axk.GetAccountX(input.ctx, toAddr)
	require.Equal(t, true, found2)
	require.Equal(t, sdk.NewInt(100000000+fee*2), input.bxk.Ak.GetAccount(input.ctx, feeAddr).GetCoins().AmountOf("cet"))
}

func TestHandleMsgSetMemoRequiredAccountNotExisted(t *testing.T) {
	input := setupTestInput()

	msg := bx.NewMsgSetTransferMemoRequired(testutil.ToAccAddress("xxx"), true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountNotActivated(t *testing.T) {
	input := setupTestInput()

	addr := testutil.ToAccAddress("myaddr")

	msg := bx.NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountOK(t *testing.T) {
	input := setupTestInput()

	addr := testutil.ToAccAddress("myaddr")
	acc := auth.NewBaseAccountWithAddress(addr)
	input.bxk.Ak.SetAccount(input.ctx, &acc)

	msg := bx.NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodeOK, result.Code)

	accX, _ := input.bxk.Axk.GetAccountX(input.ctx, addr)
	require.Equal(t, true, accX.MemoRequired)
}

func TestUnlockQueueNotAppend(t *testing.T) {
	input := setupTestInput()

	fromAddr := []byte("fromaddr")
	toAddr := []byte("toaddr")

	fromAccount := input.bxk.Ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := types2.NewAccountXWithAddress(fromAddr)

	oneCoins := dex.NewCetCoins(10100000000)
	_ = fromAccount.SetCoins(oneCoins)

	input.bxk.Ak.SetAccount(input.ctx, fromAccount)
	input.bxk.Axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := bx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 10000}
	input.handle(msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), input.bxk.Ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))
}
