package asset

import (
	"github.com/coinexchain/dex/modules/asset/internal/types"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/msgqueue"
	dex "github.com/coinexchain/dex/types"
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
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	tk  Keeper
}

func createTestInput() testInput {

	keyAsset := sdk.NewKVStoreKey(types.StoreKey)
	keyAuth := sdk.NewKVStoreKey(auth.StoreKey)
	keyAuthx := sdk.NewKVStoreKey(authx.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAsset, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAuth, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAuthx, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, false, log.NewNopLogger())
	cdc := makeTestCodec()
	types.RegisterCodec(cdc)

	// account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          {supply.Burner, supply.Minter},
		authx.ModuleName:          nil,
	}
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAuth, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, supply.DefaultCodespace, maccPerms)
	axk := authx.NewKeeper(cdc, keyAuthx, pk.Subspace(authx.DefaultParamspace), sk, ak)
	ask := NewBaseTokenKeeper(cdc, keyAsset)
	bkx := bankx.NewKeeper(pk.Subspace(bankx.DefaultParamspace), axk, bk, ak, ask, sk, msgqueue.NewProducer())
	tk := NewBaseKeeper(cdc, keyAsset, pk.Subspace(types.DefaultParamspace), bkx, sk)

	tk.SetParams(ctx, types.DefaultParams())

	initSupply := dex.NewCetCoinsE8(10000)
	sk.SetSupply(ctx, supply.NewSupply(initSupply))
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	_ = notBondedPool.SetCoins(initSupply)
	sk.SetModuleAccount(ctx, notBondedPool)

	return testInput{cdc, ctx, tk}
}

// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/staking/BaseAccount", nil)
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterConcrete(&supply.ModuleAccount{}, "test/staking/ModuleAccount", nil)
	codec.RegisterCrypto(cdc)

	return cdc
}

var _, _, testAddr = keyPubAddr()

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mockAddrList() (list []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr3, _ = sdk.AccAddressFromBech32("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q")

	list = append(list, addr1)
	list = append(list, addr2)
	list = append(list, addr3)
	return
}
