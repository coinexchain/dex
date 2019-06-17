package incentive

import (
	"github.com/coinexchain/dex/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	IncentivePoolAddr = sdk.AccAddress(crypto.AddressHash([]byte("incentive_pool")))
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	blockRewards := calcRewardsForCurrentBlock()

	collectRewardsFromPool(k, ctx, blockRewards)
}

func collectRewardsFromPool(k Keeper, ctx sdk.Context, blockRewards sdk.Coins) {
	coins, _, err := k.bankKeeper.SubtractCoins(ctx, IncentivePoolAddr, blockRewards)
	if err != nil || !coins.IsValid() {
		return
	}

	//add rewards into collected_fees for further distribution
	k.feeCollectionKeeper.AddCollectedFees(ctx, blockRewards)
}

func calcRewardsForCurrentBlock() sdk.Coins {
	//TODO: calc according incentive plan
	BlockRewardsCoins := sdk.NewCoins(sdk.NewInt64Coin(types.DefaultBondDenom, 50))
	return BlockRewardsCoins
}
