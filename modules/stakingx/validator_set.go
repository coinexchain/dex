package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slash "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	dex "github.com/coinexchain/dex/types"
)

var _ slash.StakingKeeper = Keeper{}

// forward to staking.Keeper
func (k Keeper) IterateValidators(ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.sk.IterateValidators(ctx, fn)
}
func (k Keeper) Validator(ctx sdk.Context, address sdk.ValAddress) exported.ValidatorI { // get a particular validator by operator addressL
	return k.sk.Validator(ctx, address)
}
func (k Keeper) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) exported.ValidatorI {
	return k.sk.ValidatorByConsAddr(ctx, addr)
}
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) exported.DelegationI {
	return k.sk.Delegation(ctx, addrDel, addrVal)
}
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	k.sk.Jail(ctx, consAddr)
}
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	k.sk.Unjail(ctx, consAddr)
}

func (k Keeper) MaxValidators(sdk.Context) uint16 {
	//TODO:  panic("implement me")
}

// intercept Slash(), inject CoinEx logic
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) {
	logger := ctx.Logger().With("module", "x/staking")
	oldBondedTokens := k.sk.GetPool(ctx).BondedTokens
	k.sk.Slash(ctx, consAddr, infractionHeight, power, slashFactor)

	// update pool
	pool := k.sk.GetPool(ctx)
	newBondedTokens := pool.BondedTokens
	burntTokens := oldBondedTokens.Sub(newBondedTokens)
	pool.NotBondedTokens = pool.NotBondedTokens.Add(burntTokens)
	k.sk.SetPool(ctx, pool)

	// update FeePool
	feePool := k.dk.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(dex.NewCetCoins(burntTokens.Int64()))) // TODO
	k.dk.SetFeePool(ctx, feePool)

	// TODO
	logger.Info("burnt tokens transferred from pool to fee pool!")
}
