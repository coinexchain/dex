package keepers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

func BancorInfoConsistencyInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var broken bool
		var msg string
		var count int
		keeper.Iterate(ctx, func(bi *BancorInfo) {
			if !bi.IsConsistent() {
				count++
				msg += fmt.Sprintf(" bancor Info %s consistency is broken!", bi.GetSymbol())
			}
		})
		broken = count > 0

		return sdk.FormatInvariant(types.ModuleName, "bancor info inconsistency",
			fmt.Sprintf("found %d bancor info inconsistent with initialization %s", count, msg)), broken
	}
}
