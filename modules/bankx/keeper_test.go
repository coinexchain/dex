package bankx

import (
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
	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

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
	bkxKepper := NewKeeper(subspace, authx.AccountXKeeper{}, bank.BaseKeeper{}, auth.AccountKeeper{}, auth.FeeCollectionKeeper{})

	//expect DefaultActivationFees=1
	defaultParam := DefaultParam()
	require.Equal(t, int64(100000000), defaultParam.ActivationFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { bkxKepper.SetParam(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, bkxKepper.GetParam(ctx))

}

func TestFreezeUnFreezeOK(t *testing.T) {
	input := setupTestInput()
	myaddr := testutil.ToAccAddress("myaddr")

	acc := auth.NewBaseAccountWithAddress(myaddr)
	coins := types.NewCetCoins(1000000000)
	acc.SetCoins(coins)
	input.ak.SetAccount(input.ctx, &acc)

	accx := authx.AccountX{
		Address: myaddr,
	}
	input.axk.SetAccountX(input.ctx, accx)

	freezeCoins := types.NewCetCoins(500000000)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "500000000cet", input.ak.GetAccount(input.ctx, myaddr).GetCoins().String())
	accx, _ = input.axk.GetAccountX(input.ctx, myaddr)
	require.Equal(t, "500000000cet", accx.FrozenCoins.String())

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "1000000000cet", input.ak.GetAccount(input.ctx, myaddr).GetCoins().String())
	accx, _ = input.axk.GetAccountX(input.ctx, myaddr)
	require.Equal(t, "", accx.FrozenCoins.String())
}

func TestFreezeUnFreezeInvalidAccount(t *testing.T) {
	input := setupTestInput()
	myaddr := testutil.ToAccAddress("myaddr")

	freezeCoins := types.NewCetCoins(500000000)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInvalidAddress("account doesn't exist yet"), err)

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInvalidAddress("account doesn't exist yet"), err)

}
func TestFreezeUnFreezeInsufficientCoins(t *testing.T) {
	input := setupTestInput()
	myaddr := testutil.ToAccAddress("myaddr")

	acc := auth.NewBaseAccountWithAddress(myaddr)
	coins := types.NewCetCoins(1000000000)
	acc.SetCoins(coins)
	input.ak.SetAccount(input.ctx, &acc)

	accx := authx.AccountX{
		Address: myaddr,
	}
	input.axk.SetAccountX(input.ctx, accx)

	InvalidFreezeCoins := types.NewCetCoins(5000000000)
	err := input.bxk.FreezeCoins(input.ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to freeze"), err)

	freezeCoins := types.NewCetCoins(500000000)
	err = input.bxk.FreezeCoins(input.ctx, myaddr, freezeCoins)
	require.Nil(t, err)

	err = input.bxk.UnFreezeCoins(input.ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze"), err)

}
