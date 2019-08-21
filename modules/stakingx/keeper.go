package stakingx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/stakingx/internal/types"
	dex "github.com/coinexchain/dex/types"
)

var (
	NonBondableAddressesKey = []byte("0x01")
)

type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	paramSubspace params.Subspace

	assetViewKeeper AssetViewKeeper

	sk *staking.Keeper

	dk DistributionKeeper

	ak auth.AccountKeeper

	bk ExpectBankxKeeper

	supplyKeeper ExpectSupplyKeeper

	feeCollectorName string
}

func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	paramSubspace params.Subspace, assetViewKeeper AssetViewKeeper, sk *staking.Keeper,
	dk DistributionKeeper, ak auth.AccountKeeper, bk ExpectBankxKeeper,
	supplyKeeper ExpectSupplyKeeper, feeCollectorName string) Keeper {

	return Keeper{
		key:              key,
		cdc:              cdc,
		paramSubspace:    paramSubspace.WithKeyTable(types.ParamKeyTable()),
		assetViewKeeper:  assetViewKeeper,
		sk:               sk,
		dk:               dk,
		ak:               ak,
		bk:               bk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetMinMandatoryCommissionRate(ctx sdk.Context) (rate sdk.Dec) {
	k.paramSubspace.Get(ctx, types.KeyMinMandatoryCommissionRate, &rate)
	return
}

// -----------------------------------------------------------------------------
// BondPoolStatus

func (k Keeper) CalcBondPoolStatus(ctx sdk.Context) BondPool {
	total := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(dex.CET)
	var bondPool BondPool

	bondPool.TotalSupply = total
	bondPool.BondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName).GetCoins().AmountOf(dex.CET)
	bondPool.NotBondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.NotBondedPoolName).GetCoins().AmountOf(dex.CET)
	bondPool.NonBondableTokens = calcNonBondableTokens(ctx, &k)

	bondPool.BondRatio = calcBondedRatio(&bondPool)

	return bondPool
}

func calcBondedRatio(p *BondPool) sdk.Dec {
	if p.BondedTokens.IsNegative() || p.NonBondableTokens.IsNegative() {
		return sdk.ZeroDec()
	}

	bondableTokens := p.TotalSupply.Sub(p.NonBondableTokens)
	if !bondableTokens.IsPositive() {
		return sdk.ZeroDec()
	}

	return p.BondedTokens.ToDec().QuoInt(bondableTokens)

}

func calcNonBondableTokens(ctx sdk.Context, k *Keeper) sdk.Int {
	ret := sdk.ZeroInt()
	addrs := k.getNonBondableAddresses(ctx)

	for _, addr := range addrs {
		if acc := k.ak.GetAccount(ctx, addr); acc != nil {
			amt := acc.GetCoins().AmountOf(dex.CET)
			if amt.IsPositive() {
				ret = ret.Add(amt)
			}
		}
	}

	communityPoolAmt := k.dk.GetFeePoolCommunityCoins(ctx).AmountOf(dex.CET)
	ret = ret.Add(communityPoolAmt.TruncateInt())

	return ret
}

// -----------------------------------------------------------------------------
// Non-bondable addresses

func (k Keeper) getCetOwnerAddress(ctx sdk.Context) sdk.AccAddress {
	cet := k.assetViewKeeper.GetToken(ctx, dex.CET)
	if cet == nil {
		return nil
	}
	return cet.GetOwner()
}

func (k Keeper) getAllVestingAccountAddresses(ctx sdk.Context) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0, 8)
	k.ak.IterateAccounts(ctx, func(acc auth.Account) bool {
		if vacc, ok := acc.(auth.VestingAccount); ok {
			addresses = append(addresses, vacc.GetAddress())
		}
		return false
	})
	return addresses
}

func (k Keeper) setNonBondableAddresses(ctx sdk.Context, addresses []sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinaryBare(addresses)
	if err != nil {
		panic(err)
	}
	store.Set(NonBondableAddressesKey, bz)
}

func (k Keeper) getNonBondableAddresses(ctx sdk.Context) (addresses []sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	bz := store.Get(NonBondableAddressesKey)
	if bz == nil {
		return
	}

	err := k.cdc.UnmarshalBinaryBare(bz, &addresses)
	if err != nil {
		panic(err) // TODO
	}
	return
}
