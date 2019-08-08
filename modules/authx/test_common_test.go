package authx_test

import (
	"time"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/msgqueue"

	"github.com/coinexchain/dex/modules/authx/types"

	"github.com/cosmos/cosmos-sdk/x/supply"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	dex "github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx sdk.Context
	axk authx.AccountXKeeper
	ak  auth.AccountKeeper
	sk  supply.Keeper
	cdc *codec.Codec
	tk  asset.Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	asset.RegisterCodec(cdc)

	assetKey := sdk.NewKVStoreKey("asset")
	authXKey := sdk.NewKVStoreKey("authXKey")
	authKey := sdk.NewKVStoreKey("authKey")
	keySupply := sdk.NewKVStoreKey("supply")
	skey := sdk.NewKVStoreKey("params")
	tkey := sdk.NewTransientStoreKey("transient_params")
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, "") // TODO

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	maccPerms := map[string][]string{
		types.ModuleName: nil,
		asset.ModuleName: {supply.Burner, supply.Minter},
	}

	paramsKeeper = params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	axk := authx.NewKeeper(cdc, authXKey, paramsKeeper.Subspace(authx.DefaultParamspace), sk, ak, "")
	ask := asset.NewBaseTokenKeeper(cdc, assetKey)
	bkx := bankx.NewKeeper(paramsKeeper.Subspace(bankx.DefaultParamspace), axk, bk, ak, ask, sk, msgqueue.NewProducer())
	tk := asset.NewBaseKeeper(cdc, assetKey, paramsKeeper.Subspace(asset.DefaultParamspace), bkx, sk)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id", Time: time.Unix(1560334620, 0)}, false, log.NewNopLogger())
	initSupply := dex.NewCetCoinsE8(10000)
	sk.SetSupply(ctx, supply.NewSupply(initSupply))

	return testInput{ctx: ctx, axk: axk, ak: ak, sk: sk, cdc: cdc, tk: tk}
}
