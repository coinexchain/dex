package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

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
	"github.com/coinexchain/dex/testutil"
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
	//ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	//ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))

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
	splk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{sdk.NewInt64Coin("cet", 1000000000)}}) // TODO

	initStates(ctx, sxk)
	return sxk, ctx, ak
}

func initStates(ctx sdk.Context, sxk Keeper) {
	//intialize params & states needed
	params := staking.DefaultParams()
	params.BondDenom = "cet"
	sxk.sk.SetParams(ctx, params)

	//initialize FeePool
	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(0)),
	}
	sxk.dk.SetFeePool(ctx, feePool)

	//intialize staking Pool
	pool := staking.Pool{
		NotBondedTokens: sdk.NewInt(1000e8),
		BondedTokens:    sdk.NewInt(0),
	}
	pool.String() //sxk.sk.SetPool(ctx, pool)
}

func TestSlashTokensToCommunityPool(t *testing.T) {
	sxk, ctx, ak := setUpInput()

	//initialize account for validator
	_, pk, addr := testutil.KeyPubAddr()
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(dex.NewCetCoins(1000e8))
	ak.SetAccount(ctx, &acc)

	//build createValidatorMsg
	msg := testutil.NewMsgCreateValidatorBuilder(sdk.ValAddress(addr), pk).SelfDelegation(1e8).Build()
	res := staking.NewHandler(*sxk.sk)(ctx, msg)
	staking.EndBlocker(ctx, *sxk.sk)
	bondedAmt := sdk.NewInt(1e8)
	validator := sxk.Validator(ctx, sdk.ValAddress(addr))

	//before slash
	require.True(t, res.IsOK())
	require.Equal(t, bondedAmt, validator.GetTokens())
	require.True(t, sxk.dk.GetFeePool(ctx).CommunityPool.Empty())

	//begin slash with infraction height 0
	slashfractor := sdk.NewDec(1).Quo(sdk.NewDec(20))
	sxk.Slash(ctx, sdk.ConsAddress(validator.GetConsPubKey().Address()), 0, 100, slashfractor)

	//after slash
	validator = sxk.Validator(ctx, sdk.ValAddress(addr))
	require.Equal(t, sdk.NewInt(95e6), validator.GetTokens())
	//slash tokens have been added to communityPool
	require.Equal(t, sdk.NewDec(5e6), sxk.dk.GetFeePool(ctx).CommunityPool.AmountOf("cet"))
}

func TestDelegatorSlash(t *testing.T) {
	sxk, ctx, ak := setUpInput()

	//initialize account for validator
	_, pk, addr := testutil.KeyPubAddr()
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(dex.NewCetCoins(1000e8))
	ak.SetAccount(ctx, &acc)

	//build createValidatorMsg
	msg := testutil.NewMsgCreateValidatorBuilder(sdk.ValAddress(addr), pk).SelfDelegation(1e8).Build()
	res := staking.NewHandler(*sxk.sk)(ctx, msg)
	staking.EndBlocker(ctx, *sxk.sk)
	bondedAmt := sdk.NewInt(1e8)
	validator := sxk.Validator(ctx, sdk.ValAddress(addr))

	//before slash
	require.True(t, res.IsOK())
	require.Equal(t, bondedAmt, validator.GetTokens())
	require.True(t, sxk.dk.GetFeePool(ctx).CommunityPool.Empty())

	//create new delegation at block height 2
	ctx = ctx.WithBlockHeight(2)
	_, _, addr2 := testutil.KeyPubAddr()
	acc2 := auth.NewBaseAccountWithAddress(addr2)
	acc2.SetCoins(dex.NewCetCoins(1e8))
	ak.SetAccount(ctx, &acc2)

	msgDelegation := staking.NewMsgDelegate(addr2, sdk.ValAddress(addr), dex.NewCetCoin(1e8))
	res = staking.NewHandler(*sxk.sk)(ctx, msgDelegation)

	//before slash
	validator = sxk.Validator(ctx, sdk.ValAddress(addr))
	require.True(t, res.IsOK())
	require.Equal(t, sdk.NewInt(2e8), validator.GetTokens())
	require.True(t, sxk.dk.GetFeePool(ctx).CommunityPool.Empty())

	//slash validator at block height 3
	ctx = ctx.WithBlockHeight(3)
	slashfractor := sdk.NewDec(1).Quo(sdk.NewDec(20))
	sxk.Slash(ctx, sdk.ConsAddress(validator.GetConsPubKey().Address()), 0, 200, slashfractor)

	//after slash
	validator = sxk.Validator(ctx, sdk.ValAddress(addr))
	require.Equal(t, sdk.NewInt(19e7), validator.GetTokens())

	//slash tokens have been added to communityPool
	require.Equal(t, sdk.NewDec(1e7), sxk.dk.GetFeePool(ctx).CommunityPool.AmountOf("cet"))
	delegation := sxk.Delegation(ctx, addr2, sdk.ValAddress(addr))

	//delegation at block height 2 also slashed
	tokens := validator.TokensFromShares(delegation.GetShares())
	require.Equal(t, sdk.NewDec(95e6), tokens)

	//when delegator begin undelegate, tokens he can take back are less than the amount he delegates
	msgUndelegate := staking.NewMsgUndelegate(addr2, sdk.ValAddress(addr), dex.NewCetCoin(1e8))
	res = staking.NewHandler(*sxk.sk)(ctx, msgUndelegate)
	require.Equal(t, staking.CodeInvalidDelegation, res.Code)

	msgUndelegate = staking.NewMsgUndelegate(addr2, sdk.ValAddress(addr), dex.NewCetCoin(95e6))
	res = staking.NewHandler(*sxk.sk)(ctx, msgUndelegate)
	require.True(t, res.IsOK())

	//after unbonding time, tokens are returned back to delegator's account & unbondingEntry is deleted
	ctx = ctx.WithBlockHeader(abci.Header{Time: ctx.BlockHeader().Time.Add(sxk.sk.UnbondingTime(ctx))})
	staking.EndBlocker(ctx, *sxk.sk)
	_, found := sxk.sk.GetUnbondingDelegation(ctx, addr2, sdk.ValAddress(addr))
	require.False(t, found)
	delegatorAcc := ak.GetAccount(ctx, addr2)
	require.Equal(t, dex.NewCetCoins(95e6), delegatorAcc.GetCoins())
}
