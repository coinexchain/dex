package asset

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
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

func mockWhitelist() (whitelist []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt")
	var addr2, _ = sdk.AccAddressFromBech32("cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf")
	var addr3, _ = sdk.AccAddressFromBech32("cosmos1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz")

	whitelist = append(whitelist, addr1)
	whitelist = append(whitelist, addr2)
	whitelist = append(whitelist, addr3)
	return
}

func mockAddresses() (addr []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("cosmos16cyga47yh3cv6pzemy0fjtkeqjtrjjukgngey6")
	var addr2, _ = sdk.AccAddressFromBech32("cosmos1c79cqwzah604v0pqg0h88g99p5zg08hgf0cspy")
	var addr3, _ = sdk.AccAddressFromBech32("cosmos1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk")

	addr = append(addr, addr1)
	addr = append(addr, addr2)
	addr = append(addr, addr3)
	return
}
