package incentive

import (
	"github.com/coinexchain/dex/modules/authx"
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
	"os"
	"testing"

	dex "github.com/coinexchain/dex/types"
)

type TestInput struct {
	ctx           sdk.Context
	cdc           *codec.Codec
	paramKeeper   params.Keeper
	accountKeeper auth.AccountKeeper
	keeper        Keeper
}

func SetupTestInput() TestInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)

	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	authxKey := sdk.NewKVStoreKey(authx.StoreKey)
	fckKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	incentiveKey := sdk.NewKVStoreKey(StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(incentiveKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id", Height: 10}, false, log.NewNopLogger())
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, fckKey)
	keeper := NewKeeper(cdc, incentiveKey, paramsKeeper.Subspace(DefaultParamspace), fck, bk)

	return TestInput{ctx, cdc, paramsKeeper, ak, keeper}
}

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

func TestIncentiveCoinsAddress(t *testing.T) {
	require.Equal(t, "coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97", PoolAddr.String())
}

func TestBeginBlockerInvalidCoin(t *testing.T) {

	input := SetupTestInput()
	_ = input.keeper.SetState(input.ctx, State{10})
	input.keeper.SetParam(input.ctx, DefaultParams())
	err := BeginBlocker(input.ctx, input.keeper)
	require.Equal(t, 0xa, int(err.Result().Code))
}

func TestBeginBlocker(t *testing.T) {

	input := SetupTestInput()
	_ = input.keeper.SetState(input.ctx, State{10})
	input.keeper.SetParam(input.ctx, DefaultParams())
	acc := input.accountKeeper.NewAccountWithAddress(input.ctx, PoolAddr)
	_ = acc.SetCoins(dex.NewCetCoins(10000 * 1e8))
	input.accountKeeper.SetAccount(input.ctx, acc)
	err := BeginBlocker(input.ctx, input.keeper)
	require.Equal(t, nil, err)
}
