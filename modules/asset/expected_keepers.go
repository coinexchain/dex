package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectFeeKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}
