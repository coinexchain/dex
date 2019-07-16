package asset

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"

	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
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
	tk      BaseKeeper
	handler sdk.Handler
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	RegisterCodec(cdc)
	//auth.RegisterBaseAccount(cdc)

	assetCapKey := sdk.NewKVStoreKey(StoreKey)
	authCapKey := sdk.NewKVStoreKey(auth.StoreKey)
	authxCapKey := sdk.NewKVStoreKey(authx.StoreKey)
	//fckCapKey := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	keyStaking := sdk.NewKVStoreKey(types.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(types.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authxCapKey, sdk.StoreTypeIAVL, db)
	//ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	_ = ms.LoadLatestVersion()

	var cs sdk.CodespaceType = "" // TODO
	ak := auth.NewAccountKeeper(
		cdc,
		authCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	axk := authx.NewKeeper(
		cdc,
		authxCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(authx.DefaultParamspace),
		supply.Keeper{},
		ak,
	)

	bk := bank.NewBaseKeeper(
		ak,
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(bank.DefaultParamspace),
		sdk.CodespaceRoot)
	//fck := auth.NewFeeCollectionKeeper(
	//	cdc,
	//	fckCapKey,
	//)
	ask := NewBaseTokenKeeper(
		cdc,
		assetCapKey,
	)
	bkx := bankx.NewKeeper(
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(bankx.DefaultParamspace),
		axk, bk, ak, ask,
		msgqueue.NewProducer(),
	)

	sk := staking.NewKeeper(cdc, keyStaking, tkeyStaking, nil, // TODO
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(staking.DefaultParamspace),
		types.DefaultCodespace)

	tk := NewBaseKeeper(
		cdc,
		assetCapKey,
		params.NewKeeper(cdc, keyParams, tkeyParams, cs).Subspace(DefaultParamspace),
		bkx,
		&sk)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	tk.SetParams(ctx, DefaultParams())
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
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr3, _ = sdk.AccAddressFromBech32("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q")

	whitelist = append(whitelist, addr1)
	whitelist = append(whitelist, addr2)
	whitelist = append(whitelist, addr3)
	return
}

func mockAddresses() (addr []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47")
	var addr2, _ = sdk.AccAddressFromBech32("coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g")
	var addr3, _ = sdk.AccAddressFromBech32("coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn")

	addr = append(addr, addr1)
	addr = append(addr, addr2)
	addr = append(addr, addr3)
	return
}
