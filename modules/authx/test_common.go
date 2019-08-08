package authx

import (
	"time"

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
)

type testInput struct {
	ctx sdk.Context
	axk AccountXKeeper
	ak  auth.AccountKeeper
	sk  supply.Keeper
	cdc *codec.Codec
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	authXKey := sdk.NewKVStoreKey("authXKey")
	authKey := sdk.NewKVStoreKey("authKey")
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	skey := sdk.NewKVStoreKey("params")
	tkey := sdk.NewTransientStoreKey("transient_params")
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, "") // TODO

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	maccPerms := map[string][]string{
		types.ModuleName: nil,
	}

	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	supplyKeeper := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)

	axk := NewKeeper(cdc, authXKey, paramsKeeper.Subspace(DefaultParamspace), supplyKeeper, ak, "")
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id", Time: time.Unix(1560334620, 0)}, false, log.NewNopLogger())

	return testInput{ctx: ctx, axk: axk, ak: ak, sk: supplyKeeper, cdc: cdc}
}
