package bankx

import (
	"fmt"
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

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

var myaddr = testutil.ToAccAddress("myaddr")

func defaultContext() (sdk.Context, params.Keeper) {

	cdc := codec.New()
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, skey, tkey)

	return ctx, paramsKeeper
}

func TestParamGetSet(t *testing.T) {

	ctx, paramsKeeper := defaultContext()
	subspace := paramsKeeper.Subspace(DefaultParamspace)
	bkxKepper := NewKeeper(subspace, authx.AccountXKeeper{}, bank.BaseKeeper{}, auth.AccountKeeper{}, auth.FeeCollectionKeeper{}, asset.TokenKeeper{})

	//expect DefaultActivationFees=1
	defaultParam := DefaultParams()
	require.Equal(t, int64(100000000), defaultParam.ActivationFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { bkxKepper.SetParam(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, bkxKepper.GetParam(ctx))
}

func givenAccountWith(input testInput, addr sdk.AccAddress, coinsString string) {
	coins, _ := sdk.ParseCoins(coinsString)

	acc := auth.NewBaseAccountWithAddress(addr)
	_ = acc.SetCoins(coins)
	input.ak.SetAccount(input.ctx, &acc)

	accx := authx.AccountX{
		Address: addr,
	}
	input.axk.SetAccountX(input.ctx, accx)
}

func coinsOf(input testInput, addr sdk.AccAddress) string {
	return input.ak.GetAccount(input.ctx, addr).GetCoins().String()
}

func frozenCoinsOf(input testInput, addr sdk.AccAddress) string {
	accx, _ := input.axk.GetAccountX(input.ctx, addr)
	return accx.FrozenCoins.String()
}

func TestFreezeMultiCoins(t *testing.T) {
	input := setupTestInput()

	givenAccountWith(input, myaddr, "1000000000cet,100abc")

	freezeCoins, _ := sdk.ParseCoins("300000000cet, 20abc")
	err := input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "80abc,700000000cet", coinsOf(input, myaddr))
	require.Equal(t, "20abc,300000000cet", frozenCoinsOf(input, myaddr))

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "100abc,1000000000cet", coinsOf(input, myaddr))
	require.Equal(t, "", frozenCoinsOf(input, myaddr))
}

func TestFreezeUnFreezeOK(t *testing.T) {
	input := setupTestInput()

	givenAccountWith(input, myaddr, "1000000000cet")

	freezeCoins := types.NewCetCoins(300000000)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "700000000cet", coinsOf(input, myaddr))
	require.Equal(t, "300000000cet", frozenCoinsOf(input, myaddr))

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "1000000000cet", coinsOf(input, myaddr))
	require.Equal(t, "", frozenCoinsOf(input, myaddr))
}

func TestFreezeUnFreezeInvalidAccount(t *testing.T) {
	input := setupTestInput()

	freezeCoins := types.NewCetCoins(500000000)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds;  < 500000000cet"), err)

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", myaddr)), err)
}

func TestFreezeUnFreezeInsufficientCoins(t *testing.T) {
	input := setupTestInput()

	givenAccountWith(input, myaddr, "10cet")

	InvalidFreezeCoins := types.NewCetCoins(50)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds; 10cet < 50cet"), err)

	freezeCoins := types.NewCetCoins(5)
	err = input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Nil(t, err)

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze"), err)
}
