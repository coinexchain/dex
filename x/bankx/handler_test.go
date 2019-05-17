package bankx

import (
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

func TestHandler(t *testing.T) {

	input := setupTestInput()

	fromAddr := []byte("from-address")
	toAddr := []byte("to-address")

	fromAccount := input.ak.NewAccountWithAddress(input.ctx, fromAddr)
	fromAccountX := authx.NewAccountXWithAddress(fromAddr)
	coins := sdk.NewCoins(sdk.Coin{"cet", sdk.NewInt(int64(10))})
	fromAccount.SetCoins(coins)

	input.ak.SetAccount(input.ctx, fromAccount)
	input.axk.SetAccountX(input.ctx, fromAccountX)

	msgSend := bank.MsgSend{fromAddr, toAddr, sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(int64(2))))}

	handleMsgSend(input.ctx, input.bxk, msgSend)

	require.Equal(t, sdk.NewInt(int64(8)), input.ak.GetAccount(input.ctx, fromAddr).GetCoins().AmountOf("cet"))
	require.Equal(t, sdk.NewInt(int64(1)), input.ak.GetAccount(input.ctx, toAddr).GetCoins().AmountOf("cet"))

}
