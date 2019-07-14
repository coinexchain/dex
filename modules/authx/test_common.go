package authx

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

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
	cdc *codec.Codec
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authXKey := sdk.NewKVStoreKey("authXKey")
	authKey := sdk.NewKVStoreKey("authKey")
	skey := sdk.NewKVStoreKey("params")
	tkey := sdk.NewTransientStoreKey("transient_params")
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, "") // TODO

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	axk := NewKeeper(cdc, authXKey, paramsKeeper.Subspace(bank.DefaultParamspace))
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id", Time: time.Unix(1560334620, 0)}, false, log.NewNopLogger())

	return testInput{ctx: ctx, axk: axk, ak: ak, cdc: cdc}
}
