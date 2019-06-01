package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var _, _, tAccAddr = keyPubAddr()

type testInput struct {
	cdc     *codec.Codec
	ctx     sdk.Context
	tk      TokenKeeper
	handler sdk.Handler
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	assetCapKey := sdk.NewKVStoreKey(StoreKey)
	authCapKey := sdk.NewKVStoreKey(auth.StoreKey)
	fckCapKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	_ = ms.LoadLatestVersion()

	ak := auth.NewAccountKeeper(
		cdc,
		authCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	fck := auth.NewFeeCollectionKeeper(
		cdc,
		fckCapKey,
	)
	tk := NewKeeper(
		cdc,
		assetCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(DefaultParamspace),
		ak,
		fck,
	)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	handler := NewHandler(tk)

	return testInput{cdc, ctx, tk, handler}
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
