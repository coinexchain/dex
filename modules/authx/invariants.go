package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PreTotalSupplyInvariant(k AccountXKeeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		k.PreTotalSupply(ctx)

		// TODO
		return sdk.FormatInvariant(ModuleName, "total supply",
			"ok"), false
	}
}
