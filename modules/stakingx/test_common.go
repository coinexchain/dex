package stakingx

import (
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	dex "github.com/coinexchain/dex/types"
)

func setUpInput() (Keeper, sdk.Context, auth.AccountKeeper) {
	db := dbm.NewMemDB()
	cdc := codec.New()
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	distKey := sdk.NewKVStoreKey(distribution.StoreKey)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(distKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)

	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {supply.Basic},
		authx.ModuleName:          {supply.Basic},
		distribution.ModuleName:   {supply.Basic},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		asset.ModuleName:          {supply.Minter},
	}
	splk := supply.NewKeeper(cdc, supplyKey, ak, bk, supply.DefaultCodespace, maccPerms)

	sk := staking.NewKeeper(
		cdc,
		keyStaking, tkey, splk,
		paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	dk := distribution.NewKeeper(cdc, distKey, paramsKeeper.Subspace(distribution.StoreKey), sk, splk, types.DefaultCodespace, auth.FeeCollectorName)
	sxk := NewKeeper(paramsKeeper.Subspace(DefaultParamspace), nil, &sk, dk, ak, nil, splk, auth.FeeCollectorName) // TODO

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id", Height: 1}, false, log.NewNopLogger())
	bk.SetSendEnabled(ctx, true)

	initStates(ctx, sxk, ak, splk)

	return sxk, ctx, ak
}

func initStates(ctx sdk.Context, sxk Keeper, ak auth.AccountKeeper, splk supply.Keeper) {
	//intialize params & states needed
	params := staking.DefaultParams()
	params.BondDenom = "cet"
	sxk.sk.SetParams(ctx, params)

	//initialize FeePool
	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(0)),
	}
	sxk.dk.SetFeePool(ctx, feePool)

	//initialize staking Pool
	bondedAcc := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	notBondedAcc := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	ak.SetAccount(ctx, bondedAcc)
	ak.SetAccount(ctx, notBondedAcc)

	//initialize total supply
	splk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{sdk.NewInt64Coin("cet", 10e8)}})

}
