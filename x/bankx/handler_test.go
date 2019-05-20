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

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
	"github.com/coinexchain/dex/x/authx"
)

type testInput struct {
	ctx     sdk.Context
	ak      auth.AccountKeeper
	pk      params.Keeper
	bk      bank.Keeper
	bxk     Keeper
	axk     authx.AccountXKeeper
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

	ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, fckKey)
	axk := authx.NewKeeper(cdc, authxKey)
	bxkKeeper := NewKeeper(paramsKeeper.Subspace(CodeSpaceBankx), axk, bk, ak, fck)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, DefaultParam())

	handler := NewHandler(bxkKeeper)
	return testInput{ctx: ctx, ak: ak, pk: paramsKeeper, bk: bk, bxk: bxkKeeper, axk: axk, handler: handler}
}

type testSendCases struct {
	fromAddr  string
	toAddr    string
	fromCoins sdk.Coins
	amt       sdk.Coins
}

func TestHandlerCases(t *testing.T) {

	input := setupTestInput()

	testCases := []testSendCases{
		{"fromaddr1", "toaddr1", dex.NewCetCoins(10), dex.NewCetCoins(2)},
		{"fromaddr2", "toaddr2", dex.NewCetCoins(10), dex.NewCetCoins(1)},
		{"fromaddr3", "toaddr3", dex.NewCetCoins(0), dex.NewCetCoins(2)},
	}

	var fromAccount = make([]auth.Account, len(testCases))
	var fromAccountX = make([]authx.AccountX, len(testCases))

	for i, v := range testCases {

		fromAccount[i] = input.ak.NewAccountWithAddress(input.ctx, []byte(v.fromAddr))
		fromAccountX[i] = authx.NewAccountXWithAddress([]byte(v.fromAddr))
		fromAccount[i].SetCoins(v.fromCoins)

		input.ak.SetAccount(input.ctx, fromAccount[i])
		input.axk.SetAccountX(input.ctx, fromAccountX[i])

		msgSend := bank.MsgSend{FromAddress: []byte(v.fromAddr), ToAddress: []byte(v.toAddr), Amount: v.amt}

		switch i {

		case 0:

			input.handle(msgSend)
			require.Equal(t, sdk.NewInt(int64(8)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))

			input.handle(msgSend)
			require.Equal(t, sdk.NewInt(int64(6)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(3)), input.ak.GetAccount(input.ctx, []byte(v.toAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(1)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
		case 1:
			input.handle(msgSend)
			require.Equal(t, sdk.NewInt(int64(9)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(2)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
		case 2:
			input.handle(msgSend)
			require.Equal(t, sdk.NewInt(int64(0)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(2)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
		}
	}

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
	accX := authx.NewAccountXWithAddress(addr)
	input.axk.SetAccountX(input.ctx, accX)

	msg := NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, dex.CodespaceDEX, result.Codespace)
	require.Equal(t, dex.CodeUnactivatedAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountOK(t *testing.T) {
	input := setupTestInput()

	addr := testutil.ToAccAddress("myaddr")
	accX := authx.NewAccountXWithAddress(addr)
	accX.Activated = true
	input.axk.SetAccountX(input.ctx, accX)

	accX, _ = input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t,false, accX.TransferMemoRequired)

	msg := NewMsgSetTransferMemoRequired(addr, true)
	result := input.handle(msg)
	require.Equal(t, sdk.CodeOK, result.Code)

	accX, _ = input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t,true, accX.TransferMemoRequired)
}