package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PreTotalSupplyInvariant(k AccountXKeeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		return k.PreTotalSupply(ctx)
	}
}
