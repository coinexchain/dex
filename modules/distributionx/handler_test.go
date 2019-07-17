package distributionx

import (
	"github.com/coinexchain/dex/modules/asset"
	types3 "github.com/coinexchain/dex/modules/authx/types"
	"github.com/coinexchain/dex/modules/bankx"
	types2 "github.com/coinexchain/dex/modules/distributionx/types"

	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	"testing"

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
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/types"
)

type testInput struct {
	k   Keeper
	ctx sdk.Context
	ak  auth.AccountKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	authxKey := sdk.NewKVStoreKey(authx.StoreKey)
	distrKey := sdk.NewKVStoreKey(distribution.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(distrKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	axk := authx.NewKeeper(cdc, authxKey, paramsKeeper.Subspace(types3.DefaultParamspace))
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, asset.BaseTokenKeeper{}, msgqueue.Producer{})
	distrKeeper := distribution.NewKeeper(cdc, distrKey, paramsKeeper.Subspace(distribution.DefaultParamspace), bk, staking.Keeper{}, auth.FeeCollectionKeeper{}, distribution.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	bk.SetSendEnabled(ctx, true)
	distrKeeper.SetFeePool(ctx, distribution.InitialFeePool())

	return testInput{
		ctx: ctx,
		k: Keeper{
			bxk: bxkKeeper,
			dk:  distrKeeper,
		},
		ak: ak,
	}
}
func TestDonateToCommunityPool(t *testing.T) {
	input := setupTestInput()

	addr := sdk.AccAddress([]byte("addr"))
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(types2.validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.k.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := types.NewCetCoins(1e8)
	msg := types2.NewMsgDonateToCommunityPool(addr, Coins)
	res := handleMsgDonateToCommunityPool(input.ctx, input.k, msg)

	require.True(t, res.IsOK())
	feePool = input.k.dk.GetFeePool(input.ctx)
	require.True(t, feePool.CommunityPool.AmountOf("cet").Equal(sdk.NewDec(1e8)))

	fromAcc = input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(9e8)))

}

func TestDonateToCommunityPoolFailed(t *testing.T) {
	input := setupTestInput()

	addr := sdk.AccAddress([]byte("addr"))
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(types2.validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.k.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := types.NewCetCoins(11e8)
	msg := types2.NewMsgDonateToCommunityPool(addr, Coins)
	res := handleMsgDonateToCommunityPool(input.ctx, input.k, msg)

	require.False(t, res.IsOK())
	feePool = input.k.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	fromAcc = input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

}
func TestNewHandler(t *testing.T) {
	input := setupTestInput()
	handler := NewHandler(input.k)

	msg := bankx.MsgSetMemoRequired{}
	res := handler(input.ctx, msg)

	require.Equal(t, sdk.CodeUnknownRequest, res.Code)
}
