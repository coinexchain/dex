package bankx

import (
	"github.com/coinexchain/dex/denoms"
	"github.com/coinexchain/dex/x/authx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

type testInput struct {
	ctx sdk.Context
	ak  auth.AccountKeeper
	pk  params.Keeper
	bk  bank.Keeper
	bxk Keeper
	axk authx.AccountXKeeper
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
	return testInput{ctx: ctx, ak: ak, pk: paramsKeeper, bk: bk, bxk: bxkKeeper, axk: axk}
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
		{"fromaddr1", "toaddr1", denoms.NewCetCoins(10), denoms.NewCetCoins(2)},
		{"fromaddr2", "toaddr2", denoms.NewCetCoins(10), denoms.NewCetCoins(1)},
		{"fromaddr3", "toaddr3", denoms.NewCetCoins(0), denoms.NewCetCoins(2)},
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

			handleMsgSend(input.ctx, input.bxk, msgSend)
			require.Equal(t, sdk.NewInt(int64(8)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))

			handleMsgSend(input.ctx, input.bxk, msgSend)
			require.Equal(t, sdk.NewInt(int64(6)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(3)), input.ak.GetAccount(input.ctx, []byte(v.toAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(1)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
		case 1:
			handleMsgSend(input.ctx, input.bxk, msgSend)
			require.Equal(t, sdk.NewInt(int64(9)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(2)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))
		case 2:
			handleMsgSend(input.ctx, input.bxk, msgSend)
			require.Equal(t, sdk.NewInt(int64(0)), input.ak.GetAccount(input.ctx, []byte(v.fromAddr)).GetCoins().AmountOf("cet"))
			require.Equal(t, sdk.NewInt(int64(2)), input.bxk.fck.GetCollectedFees(input.ctx).AmountOf("cet"))

		}
	}

}
