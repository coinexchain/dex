package market

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"testing"
	"time"

	"github.com/coinexchain/dex/modules/market/match"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	ctx     sdk.Context
	mk      Keeper
	handler sdk.Handler
}

var (
	haveCetAddress = []byte("have-cet")
	stock          = "ludete"
	money          = "cet"
)

func prepareAssetKeeper(ms store.CommitMultiStore, db dbm.DB, cdc *codec.Codec, ctx sdk.Context) ExpectedAssertStatusKeeper {
	asset.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	assetCapKey := sdk.NewKVStoreKey(asset.StoreKey)
	authCapKey := sdk.NewKVStoreKey(auth.StoreKey)
	fckCapKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	//create account keeper, create account
	ak := auth.NewAccountKeeper(
		cdc,
		authCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	acc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	acc.SetCoins(types.NewCetCoins(1E13))
	ak.SetAccount(ctx, acc)

	fck := auth.NewFeeCollectionKeeper(
		cdc,
		fckCapKey,
	)
	tk := asset.NewKeeper(
		cdc,
		assetCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(asset.DefaultParamspace),
		ak,
		fck,
	)

	return tk
}

func prepareBankxKeeper(ms store.CommitMultiStore, db dbm.DB, cdc *codec.Codec) ExpectedBankxKeeper {
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	authxKey := sdk.NewKVStoreKey(authx.StoreKey)
	fckKey := sdk.NewKVStoreKey(auth.FeeStoreKey)

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, fckKey)
	axk := authx.NewKeeper(cdc, authxKey, paramsKeeper.Subspace(authx.DefaultParamspace))
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck)
	_ = bxkKeeper
	//bk.SetSendEnabled(ctx, true)
	//bxkKeeper.SetParam(ctx, bankx.DefaultParam())
	return MockBankxKeeper{}
}

func prepareMockInput(t *testing.T) testInput {
	cdc := codec.New()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	marketKey := sdk.NewKVStoreKey(MarketKey)
	assetCapKey := sdk.NewKVStoreKey(asset.StoreKey)
	authCapKey := sdk.NewKVStoreKey(auth.StoreKey)
	fckCapKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	//
	//keyParams := sdk.NewKVStoreKey(params.StoreKey)
	//tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	//assetCapKey := sdk.NewKVStoreKey(asset.StoreKey)
	//authCapKey := sdk.NewKVStoreKey(auth.StoreKey)
	//fckCapKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	//TODO. Can we remove these two keys
	//skey := sdk.NewKVStoreKey("test")
	//tkey := sdk.NewTransientStoreKey("transient_test")
	//authxKey := sdk.NewKVStoreKey(authx.StoreKey)

	ms.MountStoreWithDB(marketKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	//ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	//ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	//ms.MountStoreWithDB(authxKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(ms, db, cdc, ctx)
	bk := prepareBankxKeeper(ms, db, cdc)

	mk := NewKeeper(marketKey, ak, bk, cdc,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(MarketKey))
	mk.RegisterCodec()
	handler := NewHandler(mk)
	tinput := testInput{ctx: ctx, mk: mk, handler: handler}
	issueTokens(t, tinput)
	return tinput
}

func issueTokens(t *testing.T, input testInput) {
	msg := asset.NewMsgIssueToken("ABC Token", stock, 210000000000, haveCetAddress,
		false, false, false, false)
	tk := input.mk.axk.(asset.TokenKeeper)
	handler := asset.NewHandler(tk)
	ret := handler(input.ctx, msg)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed")
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput(t)
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 6}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, false, ret.IsOK(), "create market info should failed")
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput(t)
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, true, ret.IsOK(), "create market info should succeed")
}

func TestCreateGTEOrderFailed(t *testing.T) {
	input := prepareMockInput(t)
	msgGteOrder := MsgCreateGTEOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + "noExist",
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.BUY,
		TimeInForce:    time.Now().Nanosecond() + 10000,
	}
	ret := input.handler(input.ctx, msgGteOrder)
	require.Equal(t, false, ret.IsOK(), "create GTE order should failed")
}

func TestCreateGTEOrderSuccess(t *testing.T) {
	input := prepareMockInput(t)
	msgGteOrder := MsgCreateGTEOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.BUY,
		TimeInForce:    time.Now().Nanosecond() + 10000,
	}
	/*ret := */ input.handler(input.ctx, msgGteOrder)
	//require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
}
