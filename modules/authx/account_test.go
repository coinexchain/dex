package authx

import (
	"github.com/coinexchain/dex/testutil"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"

	"github.com/stretchr/testify/require"
	"testing"
)

type testContext struct {
	ctx sdk.Context
	axk AccountXKeeper
	ak  auth.AccountKeeper
}

func setupTestCtx() testContext {
	db := dbm.NewMemDB()
	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authXKey := sdk.NewKVStoreKey("authXKey")
	authKey := sdk.NewKVStoreKey(auth.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	paramsKeeper := params.NewKeeper(cdc, skey, tkey)

	axk := NewKeeper(cdc, authXKey, paramsKeeper.Subspace(DefaultParamspace))
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return testContext{ctx: ctx, axk: axk, ak: ak}
}

func TestAccountX_GetAllUnlockedCoinsAtTheTime(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: true, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom: "eth", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	res := acc.GetAllUnlockedCoinsAtTheTime(1000)
	require.Equal(t, LockedCoins{
		LockedCoin{sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, 1000}}, res)
}

func TestAccountX_GetUnlockedCoinsAtTheTime(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: true, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom: "eth", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	res := acc.GetUnlockedCoinsAtTheTime("bch", 2000)
	require.Equal(t, LockedCoins{
		LockedCoin{sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, 1000},
		LockedCoin{sdk.Coin{Denom: "bch", Amount: sdk.NewInt(30)}, 2000},
	}, res)
}

func TestAccountX_GetAllLockedCoins(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: true, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom: "eth", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	res := acc.GetAllLockedCoins()
	require.Equal(t, coins, res)
}

func TestAccountX_GetLockedCoinsByDemon(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: true, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom: "eth", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	res := acc.GetLockedCoinsByDemon("eos")
	require.Equal(t, LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000}}, res)
}

func TestAccountX_TransferUnlockedCoins(t *testing.T) {
	ctx := setupTestCtx()
	_, pub, addr := testutil.KeyPubAddr()

	fromAccount := auth.NewBaseAccountWithAddress(addr)
	fromAccount.SetPubKey(pub)
	oneCoins := sdk.Coins{sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}}
	fromAccount.SetCoins(oneCoins)

	ctx.ak.SetAccount(ctx.ctx, &fromAccount)

	var acc = AccountX{Address: addr, Activated: true, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom: "eth", Amount: sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom: "eos", Amount: sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	ctx.axk.SetAccountX(ctx.ctx, acc)

	acc.TransferUnlockedCoins(1000, ctx.ctx, ctx.axk, ctx.ak)
	require.Equal(t, "eth", acc.LockedCoins[0].Coin.Denom)
	require.Equal(t, "eos", acc.LockedCoins[1].Coin.Denom)

	require.Equal(t, sdk.NewInt(40), ctx.ak.GetAccount(ctx.ctx, addr).GetCoins().AmountOf("bch"))
}

func TestAccountX_SetActivated(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: false, MemoRequired: false}
	acc.SetActivated(true)
	require.Equal(t, true, acc.IsActivated())
}

func TestAccountX_IsActivated(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: true, MemoRequired: false}
	res := acc.IsActivated()
	require.Equal(t, true, res)
}

func TestAccountX_AddLockedCoins(t *testing.T) {
	var acc = AccountX{Address: []byte("123"), Activated: false, MemoRequired: false}
	acc.AddLockedCoins(LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "bch", Amount: sdk.NewInt(10)}, UnlockTime: 1000}})
	require.Equal(t, "bch", acc.GetLockedCoinsByDemon("bch")[0].Coin.Denom)
	require.Equal(t, sdk.NewInt(10), acc.GetLockedCoinsByDemon("bch")[0].Coin.Amount)
}
