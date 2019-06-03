package market

import (
	"testing"
	"time"

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

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx     sdk.Context
	mk      Keeper
	handler sdk.Handler
}

var (
	haveCetAddress    = []byte("have-cet")
	notHaveCetAddress = []byte("no-have-cet")
	stock             = "usdt"
	money             = "eos"
)

type storeKeys struct {
	assetCapKey *sdk.KVStoreKey
	authCapKey  *sdk.KVStoreKey
	fckCapKey   *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey
	marketKey   *sdk.KVStoreKey
	authxKey    *sdk.KVStoreKey
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

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) ExpectedBankxKeeper {

	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams)
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, keys.fckCapKey)
	axk := authx.NewKeeper(cdc, keys.authxKey, paramsKeeper.Subspace(authx.DefaultParamspace))
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck)
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, bankx.DefaultParams())

	return bxkKeeper
}

func prepareMockInput(t *testing.T) testInput {
	cdc := codec.New()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := storeKeys{}
	keys.marketKey = sdk.NewKVStoreKey(StoreKey)
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.fckCapKey = sdk.NewKVStoreKey(auth.FeeStoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(t, keys, cdc, ctx)
	bk := prepareBankxKeeper(keys, cdc, ctx)

	mk := NewKeeper(keys.marketKey, ak, bk, cdc, params.NewKeeper(
		cdc, keys.keyParams, keys.tkeyParams).Subspace(StoreKey))
	mk.RegisterCodec()
	return testInput{ctx: ctx, mk: mk, handler: NewHandler(mk)}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput(t)
	// failed by price precision
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 6}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create market info should failed")

	// failed by token
	msgMarketInfo.Money = "btc"
	msgMarketInfo.PricePrecision = 8
	ret = input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, CodeInvalidToken, ret.Code, "create market info should failed")

	// failed by coins
	msgMarketInfo.Money = money
	msgMarketInfo.Creator = notHaveCetAddress
	ret = input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, CodeInvalidTokenIssuer, ret.Code, "create market info should failed")

}

func createMarket(input testInput) sdk.Result {
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	return input.handler(input.ctx, msgMarketInfo)
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput(t)
	ret := createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market info should succeed")
}

func TestCreateGTEOrderFailed(t *testing.T) {
	input := prepareMockInput(t)
	msgGteOrder := MsgCreateOrder{
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
	createMarket(input)
	ret := input.handler(input.ctx, msgGteOrder)
	require.Equal(t, false, ret.IsOK(), "create GTE order should failed")
}

func TestCreateGTEOrderSuccess(t *testing.T) {
	input := prepareMockInput(t)
	msgGteOrder := MsgCreateOrder{
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
	ret := createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")
	ret = input.handler(input.ctx, msgGteOrder)
	require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
}
