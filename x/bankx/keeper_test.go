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

func defaultContext(key sdk.StoreKey, tkey sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestParamGetSet(t *testing.T) {
	cdc := codec.New()
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	ctx := defaultContext(skey, tkey)

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	subspace := paramsKeeper.Subspace(DefaultParamSpace)
	bkxKepper := NewKeeper(subspace, authx.AccountXKeeper{}, bank.BaseKeeper{}, auth.FeeCollectionKeeper{})

	//expect DefaultActivatedFees=1
	defaultParam := DefaultParam()
	require.Equal(t, int64(1), defaultParam.ActivatedFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { bkxKepper.SetParam(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, bkxKepper.GetParam(ctx))

}
