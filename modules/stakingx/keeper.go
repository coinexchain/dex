package stakingx

import (
	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "stakingx"

	StoreKey = ModuleName

	// DefaultParamspace defines the default stakingx module parameter subspace
	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for stakingx
	QuerierRoute = ModuleName
)

type Keeper struct {
	paramSubspace params.Subspace

	assetViewKeeper AssetViewKeeper

	sk *staking.Keeper

	dk DistributionKeeper

	ak auth.AccountKeeper

	bk ExpectBankxKeeper

	supplyKeeper ExpectSupplyKeeper

	feeCollectorName string
}

func NewKeeper(paramSubspace params.Subspace, assetViewKeeper AssetViewKeeper, sk *staking.Keeper,
	dk DistributionKeeper, ak auth.AccountKeeper, bk ExpectBankxKeeper,
	supplyKeeper ExpectSupplyKeeper, feeCollectorName string) Keeper {

	return Keeper{
		paramSubspace:    paramSubspace.WithKeyTable(ParamKeyTable()),
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
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetNonBondableAddresses(ctx sdk.Context) (addrs []sdk.AccAddress) {
	k.paramSubspace.Get(ctx, KeyNonBondableAddresses, &addrs)
	return
}

func (k Keeper) GetMinMandatoryCommissionRate(ctx sdk.Context) (rate sdk.Dec) {
	k.paramSubspace.Get(ctx, KeyMinMandatoryCommissionRate, &rate)
	return
}

func (k Keeper) CalcBondPoolStatus(ctx sdk.Context) BondPool {

	total := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(types.CET)
	var bondPool BondPool

	bondPool.TotalSupply = total
	bondPool.BondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName).GetCoins().AmountOf(types.CET)
	bondPool.NotBondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.NotBondedPoolName).GetCoins().AmountOf(types.CET)
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
	addrs := k.GetNonBondableAddresses(ctx)

	for _, addr := range addrs {
		if acc := k.ak.GetAccount(ctx, addr); acc != nil {
			amt := acc.GetCoins().AmountOf(types.CET)
			if amt.IsPositive() {
				ret = ret.Add(amt)
			}
		}
	}

	communityPoolAmt := k.dk.GetFeePoolCommunityCoins(ctx).AmountOf(types.CET)
	ret = ret.Add(communityPoolAmt.TruncateInt())

	return ret
}
