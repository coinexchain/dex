package incentive

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/msgqueue"
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
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"testing"

	dex "github.com/coinexchain/dex/types"
)

type fakeAssetStatusKeeper struct{}

func (k fakeAssetStatusKeeper) IsTokenForbidden(ctx sdk.Context, symbol string) bool {
	return false
}
func (k fakeAssetStatusKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool {
	return false
}
func (k fakeAssetStatusKeeper) UpdateTokenSendLock(ctx sdk.Context, symbol string, amount sdk.Int, lock bool) sdk.Error {
	return nil
}

func defaultContext() (sdk.Context, *codec.Codec, Keeper, auth.AccountKeeper) {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	bank.RegisterCodec(cdc)
	bankx.RegisterCodec(cdc)
	RegisterCodec(cdc)

	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")

	keySupply := sdk.NewKVStoreKey("supply")
	keyAuth := sdk.NewKVStoreKey("auth")
	keyAuthX := sdk.NewKVStoreKey("authx")
	keyBank := sdk.NewKVStoreKey("bank")
	keyBankx := sdk.NewKVStoreKey("bankx")
	KeyIncentive := sdk.NewKVStoreKey(StoreKey)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyAuthX, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyBankx, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(KeyIncentive, sdk.StoreTypeIAVL, db)
	_ = cms.LoadLatestVersion()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		authx.ModuleName:          nil,
		bank.ModuleName:           nil,
		bankx.ModuleName:          nil,
		ModuleName:                nil,
	}

	ask := fakeAssetStatusKeeper{}

	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAuth, paramsKeeper.Subspace("auth"), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace("bank"), "bank", map[string]bool{})
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	axk := authx.NewKeeper(cdc, keyAuthX, paramsKeeper.Subspace("authx"), sk, ak, "")
	_ = bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, ask, sk, msgqueue.NewProducer())
	keeper := NewKeeper(cdc, KeyIncentive, paramsKeeper.Subspace(ModuleName), bk, sk, "FeeCollector")
	return ctx, cdc, keeper, ak
}

type TestInput struct {
	ctx    sdk.Context
	cdc    *codec.Codec
	keeper Keeper
	ak     auth.AccountKeeper
}

func SetupTestInput() TestInput {
	ctx, cdc, Keep, ak := defaultContext()
	return TestInput{ctx: ctx, cdc: cdc, keeper: Keep, ak: ak}
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
	acc := input.ak.NewAccountWithAddress(input.ctx, PoolAddr)
	_ = acc.SetCoins(dex.NewCetCoins(10000 * 1e8))
	input.ak.SetAccount(input.ctx, acc)
	err := BeginBlocker(input.ctx, input.keeper)
	require.Equal(t, nil, err)
}
