package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//type testContext struct {
//	ctx sdk.Context
//	axk AccountXKeeper
//	ak  auth.AccountKeeper
//}

//func setupTestCtx() testContext {
//	db := dbm.NewMemDB()
//	cdc := codec.New()
//	RegisterCodec(cdc)
//	auth.RegisterCodec(cdc)
//	codec.RegisterCrypto(cdc)
//
//	authXKey := sdk.NewKVStoreKey("authXKey")
//	authKey := sdk.NewKVStoreKey(auth.StoreKey)
//
//	ms := store.NewCommitMultiStore(db)
//	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
//	ms.LoadLatestVersion()
//
//	skey := sdk.NewKVStoreKey("test")
//	tkey := sdk.NewTransientStoreKey("transient_test")
//	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
//
//	axk := NewKeeper(cdc, authXKey, paramsKeeper.Subspace(DefaultParamspace))
//	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
//	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
//
//	return testContext{ctx: ctx, axk: axk, ak: ak}
//}

func TestAccountX_GetAllUnlockedCoinsAtTheTime(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	coins := LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("eth", sdk.NewInt(30), 2000),
		NewLockedCoin("eos", sdk.NewInt(40), 3000),
	}
	acc.LockedCoins = coins
	res := acc.GetAllUnlockedCoinsAtTheTime(1000)
	require.Equal(t, LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000)}, res)
}

func TestAccountX_GetUnlockedCoinsAtTheTime(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	coins := LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("eth", sdk.NewInt(30), 2000),
		NewLockedCoin("bch", sdk.NewInt(30), 2000),
		NewLockedCoin("eos", sdk.NewInt(40), 3000),
	}
	acc.LockedCoins = coins
	res := acc.GetUnlockedCoinsAtTheTime("bch", 2000)
	require.Equal(t, LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("bch", sdk.NewInt(30), 2000),
	}, res)
}

func TestAccountX_GetAllLockedCoins(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	coins := LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("eth", sdk.NewInt(30), 2000),
		NewLockedCoin("eos", sdk.NewInt(40), 3000),
	}
	acc.LockedCoins = coins
	res := acc.GetAllLockedCoins()
	require.Equal(t, coins, res)
}

func TestAccountX_GetLockedCoinsByDemon(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	coins := LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("eth", sdk.NewInt(30), 2000),
		NewLockedCoin("eos", sdk.NewInt(40), 3000),
	}
	acc.LockedCoins = coins
	res := acc.GetLockedCoinsByDemon("eos")
	require.Equal(t, LockedCoins{
		NewLockedCoin("eos", sdk.NewInt(40), 3000)}, res)
}

func TestAccountX_AddLockedCoins(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	acc.AddLockedCoins(LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(10), 1000)})
	require.Equal(t, "bch", acc.GetLockedCoinsByDemon("bch")[0].Coin.Denom)
	require.Equal(t, sdk.NewInt(10), acc.GetLockedCoinsByDemon("bch")[0].Coin.Amount)
}

func TestAccountX_GetAllCoins(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), MemoRequired: false}
	coins := LockedCoins{
		NewLockedCoin("bch", sdk.NewInt(20), 1000),
		NewLockedCoin("eth", sdk.NewInt(30), 2000),
		NewLockedCoin("eos", sdk.NewInt(40), 3000),
	}
	acc.LockedCoins = coins
	acc.FrozenCoins = sdk.NewCoins(
		sdk.NewCoin("bch", sdk.NewInt(50)),
		sdk.NewCoin("eth", sdk.NewInt(10)),
	)

	res := acc.GetAllCoins()
	expected := sdk.NewCoins(
		sdk.NewCoin("bch", sdk.NewInt(70)),
		sdk.NewCoin("eth", sdk.NewInt(40)),
		sdk.NewCoin("eos", sdk.NewInt(40)),
	)

	require.Equal(t, expected, res)
}
