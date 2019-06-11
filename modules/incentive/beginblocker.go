package incentive

import (
	"github.com/coinexchain/dex/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BlockRewards = 50

var (
	IncentiveCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("incentive_pool")))
	BlockRewardsCoins     = sdk.NewCoins(sdk.NewInt64Coin(types.DefaultBondDenom, BlockRewards))
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	coins, _, err := k.bankKeeper.SubtractCoins(ctx, IncentiveCoinsAccAddr, BlockRewardsCoins)
	if err != nil || !coins.IsValid() {
		return
	}

	incentiveCoins := sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(int64(BlockRewards))))
	//TODO:sub the corresponding coins from genesis account

	//add these coins to collectedfees
	k.feeCollectionKeeper.AddCollectedFees(ctx, incentiveCoins)
}
