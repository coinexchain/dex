package market

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	//abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	//"github.com/tendermint/tendermint/libs/log"
	"testing"
)

type testBankxInput struct {
	ak  auth.AccountKeeper
	bxk bankx.Keeper
}

type testAssetInput struct {
}

type testInput struct {
	ctx sdk.Context
	bni testBankxInput
	asi testAssetInput
}

var (
	haveCoinAddress = "have-coin"
)

func prepareBankx(db dbm.DB, ctx sdk.Context) testBankxInput {
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
	axk := authx.NewKeeper(cdc, authxKey, paramsKeeper.Subspace(authx.StoreKey))
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck)
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, bankx.DefaultParam())

	input := testBankxInput{ak: ak, bxk: bxkKeeper}
	//input.ak.NewAccountWithAddress(ctx, haveCoinAddress)
	return input
}

func prepareAsset(db dbm.DB, ctx sdk.Context) {

}

//func setupTestInput(t *testing.T) testInput {
//db := dbm.NewMemDB()
//cdc := codec.New()
//RegisterCodec(cdc)
//marketKey := sdk.NewKVStoreKey(MarketKey)
//
//ms := store.NewCommitMultiStore(db)
//ms.MountStoreWithDB(marketKey, sdk.StoreTypeIAVL, db)
//
//mk := NewKeeper(marketKey, MockAssertKeeper{}, MockBankxKeeper{})
//ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
//input := testInput{ctx, mk}
//
////create account
//addr := sdk.AccAddress([]byte("some-address"))
//acc := authx.NewAccountXWithAddress(addr)
//require.Equal(t, addr, acc.Address)
//input.

//return input
//}

func TestMarketInfoSetFailed(t *testing.T) {
	//input := setupTestInput(t)

}

func TestMarketInfoSetSuccess(t *testing.T) {

}

func TestCreateGTEOrderFailed(t *testing.T) {

}

func TestCreateGTEOrderSuccess(t *testing.T) {

}
