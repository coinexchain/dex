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
	stock          = "usdt"
	money          = "eos"
)

type storeKeys struct {
	assetCapKey *sdk.KVStoreKey
	authCapKey  *sdk.KVStoreKey
	fckCapKey   *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey
	marketKey   *sdk.KVStoreKey
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context) ExpectedAssertStatusKeeper {
	asset.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	//create auth, asset keeper
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, params.NewKeeper(cdc, keys.keyParams,
		keys.tkeyParams).Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	fck := auth.NewFeeCollectionKeeper(cdc, keys.fckCapKey)
	tk := asset.NewKeeper(cdc, keys.assetCapKey, params.NewKeeper(cdc, keys.keyParams,
		keys.tkeyParams).Subspace(asset.DefaultParamspace), ak, fck)
	tk.SetParams(ctx, asset.DefaultParams())

	// create an account by auth keeper
	acc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	acc.SetCoins(types.NewCetCoins(1E13))
	ak.SetAccount(ctx, acc)

	// issue tokens
	msgStock := asset.NewMsgIssueToken(stock, stock, 210000000000, haveCetAddress,
		false, false, false, false)
	msgMoney := asset.NewMsgIssueToken(money, money, 210000000000, haveCetAddress,
		false, false, false, false)
	handler := asset.NewHandler(tk)
	ret := handler(ctx, msgStock)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)
	ret = handler(ctx, msgMoney)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)

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

	keys := storeKeys{}
	keys.marketKey = sdk.NewKVStoreKey(MarketKey)
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.fckCapKey = sdk.NewKVStoreKey(auth.FeeStoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	//TODO. Can we remove these two keys
	//skey := sdk.NewKVStoreKey("test")
	//tkey := sdk.NewTransientStoreKey("transient_test")
	//authxKey := sdk.NewKVStoreKey(authx.StoreKey)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(t, keys, cdc, ctx)
	bk := prepareBankxKeeper(ms, db, cdc)

	mk := NewKeeper(keys.marketKey, ak, bk, cdc, params.NewKeeper(
		cdc, keys.keyParams, keys.tkeyParams).Subspace(MarketKey))
	mk.RegisterCodec()
	return testInput{ctx: ctx, mk: mk, handler: NewHandler(mk)}
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
