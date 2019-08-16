package distributionx

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	types2 "github.com/coinexchain/dex/modules/distributionx/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

var validCoins = dex.NewCetCoins(10e8)

type testInput struct {
	k   Keeper
	ctx sdk.Context
	ak  auth.AccountKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	authxKey := sdk.NewKVStoreKey(authx.StoreKey)
	distrKey := sdk.NewKVStoreKey(distribution.StoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(distrKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot, map[string]bool{})

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		authx.ModuleName:          nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		asset.ModuleName:          {supply.Minter},
	}
	sk := supply.NewKeeper(cdc, supplyKey, ak, bk, maccPerms)
	//ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	//ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))

	axk := authx.NewKeeper(cdc, authxKey, paramsKeeper.Subspace(authx.DefaultParamspace), sk, ak, "")
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, asset.BaseTokenKeeper{}, sk, msgqueue.NewProducer())
	distrKeeper := distribution.NewKeeper(cdc, distrKey, paramsKeeper.Subspace(distribution.DefaultParamspace),
		staking.Keeper{}, sk, distribution.DefaultCodespace, auth.FeeCollectorName, map[string]bool{})

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
	acc.SetCoins(validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.k.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := dex.NewCetCoins(1e8)
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
	acc.SetCoins(validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.k.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := dex.NewCetCoins(11e8)
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
