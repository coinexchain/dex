package incentive

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	dex "github.com/coinexchain/dex/types"
)

var (
	PoolAddr = sdk.AccAddress(crypto.AddressHash([]byte("incentive_pool")))
)

func BeginBlocker(ctx sdk.Context, k keepers.Keeper) {
	blockRewards, planned := calcRewardsInPlans(ctx, k)
	if !planned {
		blockRewards = calDefaultRewards(ctx, k)
	}

	if k.HasCoins(ctx, PoolAddr, blockRewards) {
		if blockRewards != nil {
			if err := collectRewardsFromPool(k, ctx, blockRewards); err != nil {
				panic(err)
			}
		}
	}

}

func collectRewardsFromPool(k keepers.Keeper, ctx sdk.Context, blockRewards sdk.Coins) sdk.Error {
	//transfer rewards into collected_fees for further distribution
	err := k.SendCoinsFromAccountToModule(ctx, PoolAddr, auth.FeeCollectorName, blockRewards)
	if err != nil {
		return err
	}
	return nil
}

func calcRewardsInPlans(ctx sdk.Context, k keepers.Keeper) (sdk.Coins, bool) {
	var rewardAmount int64
	var heightInplan bool
	height := ctx.BlockHeader().Height
	adjustmentHeight := k.GetState(ctx).HeightAdjustment
	height = height + adjustmentHeight

	plans := k.GetParams(ctx).Plans

	for _, plan := range plans {
		if height > plan.StartHeight && height <= plan.EndHeight {
			rewardAmount = rewardAmount + plan.RewardPerBlock
			heightInplan = true
			break
		}
	}
	blockRewardsCoins := sdk.NewCoins(sdk.NewInt64Coin(dex.DefaultBondDenom, rewardAmount))
	return blockRewardsCoins, heightInplan
}
func calDefaultRewards(ctx sdk.Context, k keepers.Keeper) (res sdk.Coins) {
	defaultIncentive := k.GetParams(ctx).DefaultRewardPerBlock
	res = sdk.NewCoins(sdk.NewInt64Coin(dex.CET, defaultIncentive))
	return
}
