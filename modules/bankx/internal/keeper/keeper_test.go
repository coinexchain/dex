package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"

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

	"github.com/coinexchain/dex/modules/authx"
	ax "github.com/coinexchain/dex/modules/authx/types"
	bx "github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

type fakeAssetStatusKeeper struct{}

func (k fakeAssetStatusKeeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	return false
}
func (k fakeAssetStatusKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	return false
}

var myaddr = testutil.ToAccAddress("myaddr")

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
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace("bank"), "bank")
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, supply.DefaultCodespace, maccPerms)
	axk := authx.NewKeeper(cdc, keyAuthX, paramsKeeper.Subspace("authx"), sk, ak)
	bxK := NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, ask, sk, msgqueue.NewProducer())

	return ctx, cdc, bxK
}

func TestParamGetSet(t *testing.T) {
	ctx, _, keeper := defaultContext()

	//expect DefaultActivationFees=1
	defaultParam := bx.DefaultParams()
	require.Equal(t, int64(100000000), defaultParam.ActivationFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { keeper.SetParam(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, keeper.GetParam(ctx))
}

func givenAccountWith(ctx sdk.Context, keeper Keeper, addr sdk.AccAddress, coinsString string) {
	coins, _ := sdk.ParseCoins(coinsString)

	acc := auth.NewBaseAccountWithAddress(addr)
	_ = acc.SetCoins(coins)
	keeper.Ak.SetAccount(ctx, &acc)

	accX := ax.AccountX{
		Address: addr,
	}
	keeper.Axk.SetAccountX(ctx, accX)
}

func coinsOf(ctx sdk.Context, keeper Keeper, addr sdk.AccAddress) string {
	return keeper.Ak.GetAccount(ctx, addr).GetCoins().String()
}

func frozenCoinsOf(ctx sdk.Context, keeper Keeper, addr sdk.AccAddress) string {
	accX, _ := keeper.Axk.GetAccountX(ctx, addr)
	return accX.FrozenCoins.String()
}

func TestFreezeMultiCoins(t *testing.T) {
	ctx, _, keeper := defaultContext()

	givenAccountWith(ctx, keeper, myaddr, "1000000000cet,100abc")

	freezeCoins, _ := sdk.ParseCoins("300000000cet, 20abc")
	err := keeper.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "80abc,700000000cet", coinsOf(ctx, keeper, myaddr))
	require.Equal(t, "20abc,300000000cet", frozenCoinsOf(ctx, keeper, myaddr))

	err = keeper.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "100abc,1000000000cet", coinsOf(ctx, keeper, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, keeper, myaddr))
}

func TestFreezeUnFreezeOK(t *testing.T) {

	ctx, _, keeper := defaultContext()

	givenAccountWith(ctx, keeper, myaddr, "1000000000cet")

	freezeCoins := types.NewCetCoins(300000000)
	err := keeper.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "700000000cet", coinsOf(ctx, keeper, myaddr))
	require.Equal(t, "300000000cet", frozenCoinsOf(ctx, keeper, myaddr))

	err = keeper.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "1000000000cet", coinsOf(ctx, keeper, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, keeper, myaddr))
}

func TestFreezeUnFreezeInvalidAccount(t *testing.T) {

	ctx, _, keeper := defaultContext()

	freezeCoins := types.NewCetCoins(500000000)
	err := keeper.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds;  < 500000000cet"), err)

	err = keeper.UnFreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", myaddr)), err)
}

func TestFreezeUnFreezeInsufficientCoins(t *testing.T) {
	ctx, _, keeper := defaultContext()

	givenAccountWith(ctx, keeper, myaddr, "10cet")

	InvalidFreezeCoins := types.NewCetCoins(50)
	err := keeper.FreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds; 10cet < 50cet"), err)

	freezeCoins := types.NewCetCoins(5)
	err = keeper.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Nil(t, err)

	err = keeper.UnFreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze"), err)
}

func TestGetTotalCoins(t *testing.T) {
	ctx, _, keeper := defaultContext()
	givenAccountWith(ctx, keeper, myaddr, "100cet, 20bch, 30btc")

	lockedCoins := ax.LockedCoins{
		ax.NewLockedCoin("bch", sdk.NewInt(20), 1000),
		ax.NewLockedCoin("eth", sdk.NewInt(30), 2000),
	}

	frozenCoins := sdk.NewCoins(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(50)},
		sdk.Coin{Denom: "eth", Amount: sdk.NewInt(10)},
	)

	accX := ax.AccountX{
		Address:     myaddr,
		LockedCoins: lockedCoins,
		FrozenCoins: frozenCoins,
	}

	keeper.Axk.SetAccountX(ctx, accX)

	expected := sdk.NewCoins(
		sdk.Coin{Denom: "bch", Amount: sdk.NewInt(40)},
		sdk.Coin{Denom: "btc", Amount: sdk.NewInt(80)},
		sdk.Coin{Denom: "cet", Amount: sdk.NewInt(100)},
		sdk.Coin{Denom: "eth", Amount: sdk.NewInt(40)},
	)
	expected = expected.Sort()
	coins := keeper.GetTotalCoins(ctx, myaddr)

	require.Equal(t, expected, coins)
}

func TestKeeper_TotalAmountOfCoin(t *testing.T) {

	ctx, _, keeper := defaultContext()
	amount := keeper.TotalAmountOfCoin(ctx, "cet")
	require.Equal(t, int64(0), amount.Int64())

	givenAccountWith(ctx, keeper, myaddr, "100cet")

	lockedCoins := ax.LockedCoins{
		ax.NewLockedCoin("cet", sdk.NewInt(100), 1000),
	}
	frozenCoins := sdk.NewCoins(sdk.Coin{Denom: "cet", Amount: sdk.NewInt(100)})

	accX := ax.AccountX{
		Address:     myaddr,
		LockedCoins: lockedCoins,
		FrozenCoins: frozenCoins,
	}
	keeper.Axk.SetAccountX(ctx, accX)
	amount = keeper.TotalAmountOfCoin(ctx, "cet")
	require.Equal(t, int64(300), amount.Int64())
}

func TestKeeper_AddCoins(t *testing.T) {
	ctx, _, keeper := defaultContext()
	coins := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(10)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(20)},
	)

	coins2 := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(5)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(10)},
	)

	err := keeper.AddCoins(ctx, myaddr, coins)
	require.Equal(t, nil, err)
	err = keeper.SubtractCoins(ctx, myaddr, coins2)
	require.Equal(t, nil, err)
	cs := keeper.GetTotalCoins(ctx, myaddr)
	require.Equal(t, coins2, cs)

	coins3 := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(15)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(10)},
	)
	err = keeper.SubtractCoins(ctx, myaddr, coins3)
	require.Error(t, err)
}

func TestKeeper_SendCoins(t *testing.T) {
	ctx, _, keeper := defaultContext()
	coins := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(10)},
	)
	addr2 := testutil.ToAccAddress("addr2")
	_ = keeper.AddCoins(ctx, myaddr, coins)
	exist := keeper.HasCoins(ctx, myaddr, coins)
	assert.True(t, exist)
	err := keeper.SendCoins(ctx, myaddr, addr2, coins)
	require.Equal(t, nil, err)
	cs := keeper.GetTotalCoins(ctx, addr2)
	require.Equal(t, coins, cs)
}
