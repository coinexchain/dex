package stakingx

import (
	"fmt"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SupplyCETInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		totalAmount := k.bk.TotalAmountOfCoin(ctx, types.CET)

		token := k.assetViewKeeper.GetToken(ctx, types.DefaultBondDenom)
		if totalAmount.Int64() != token.GetTotalSupply() {
			return fmt.Errorf("the cet total amount [ %d ]is inconsistent with the actual amount [ %d ]",
				token.GetTotalSupply(), totalAmount.Int64())
		}

		return nil
	}
}
