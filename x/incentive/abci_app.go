package incentive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BlockRewards = 50

func BeginBlocker(ctx sdk.Context, k Keeper) {

	//TODO: check whether there is still block incentive left

	incentiveCoins := sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(int64(BlockRewards))))

	//TODO:sub the corresponding coins from genesis account

	//add these coins to collectedfees
	k.feeCollectionKeeper.AddCollectedFees(ctx, incentiveCoins)

}
