package incentive

import (
	"github.com/coinexchain/dex/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	PoolAddr = sdk.AccAddress(crypto.AddressHash([]byte("incentive_pool")))
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	blockRewards := calcRewardsForCurrentBlock(ctx, k)

	collectRewardsFromPool(k, ctx, blockRewards)
}

func collectRewardsFromPool(k Keeper, ctx sdk.Context, blockRewards sdk.Coins) {
	coins, _, err := k.bankKeeper.SubtractCoins(ctx, PoolAddr, blockRewards)
	if err != nil || !coins.IsValid() {
		return
	}

	//add rewards into collected_fees for further distribution
	k.feeCollectionKeeper.AddCollectedFees(ctx, blockRewards)
}

func calcRewardsForCurrentBlock(ctx sdk.Context, k Keeper) sdk.Coins {

	var RewardAmount int64 = 0
	height := ctx.BlockHeader().Height
	adjustmentHeight := k.GetState(ctx).HeightAdjustment
	height = height + adjustmentHeight
	plans := k.GetParam(ctx).Incentive.Plans
	for _, plan := range plans {
		if height > plan.StartHeight && height <= plan.EndHeight {
			RewardAmount = RewardAmount + plan.RewardPerBlock
		}
	}
	BlockRewardsCoins := sdk.NewCoins(sdk.NewInt64Coin(types.DefaultBondDenom, 0))
	return BlockRewardsCoins
}
