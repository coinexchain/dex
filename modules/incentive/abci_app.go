package incentive

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BlockRewards = 50

var (
	IncentiveCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("IncentiveCoins")))
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	//TODO: check whether there is still block incentive left

	incentiveCoins := sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(int64(BlockRewards))))

	//TODO:sub the corresponding coins from genesis account

	//add these coins to collectedfees
	k.feeCollectionKeeper.AddCollectedFees(ctx, incentiveCoins)

}
