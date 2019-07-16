package incentive

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/types"
)

var (
	PoolAddr = sdk.AccAddress(crypto.AddressHash([]byte("incentive_pool")))
)

func BeginBlocker(ctx sdk.Context, k Keeper) sdk.Error {
	blockRewards := calcRewardsForCurrentBlock(ctx, k)
	err := collectRewardsFromPool(k, ctx, blockRewards)
	return err
}

func collectRewardsFromPool(k Keeper, ctx sdk.Context, blockRewards sdk.Coins) sdk.Error {

	//transfer rewards into collected_fees for further distribution
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, PoolAddr, auth.FeeCollectorName, blockRewards)
	if err != nil {
		return err
	}
	return nil
}

func calcRewardsForCurrentBlock(ctx sdk.Context, k Keeper) sdk.Coins {

	var rewardAmount int64
	height := ctx.BlockHeader().Height
	adjustmentHeight := k.GetState(ctx).HeightAdjustment
	height = height + adjustmentHeight
	plans := k.GetParam(ctx).Plans
	for _, plan := range plans {
		if height > plan.StartHeight && height <= plan.EndHeight {
			rewardAmount = rewardAmount + plan.RewardPerBlock
		}
	}
	blockRewardsCoins := sdk.NewCoins(sdk.NewInt64Coin(types.DefaultBondDenom, rewardAmount))
	return blockRewardsCoins
}
