package crisisx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/types"
)

const (
	TotalCetAmount int64 = 588800000000000000
)

func Register(k Keeper) {
	k.ck.RegisterRoute(ModuleName, "cet-invariant", cetTotalInvariant(k))
}

func cetTotalInvariant(k Keeper) sdk.Invariant {

	return func(ctx sdk.Context) error {

		// Get all amounts based on the account system
		basedAccountTotalAmount := k.bk.TotalAmountOfCoin(ctx, types.CET)
		feeAmount := k.feek.GetCollectedFees(ctx).AmountOf(types.CET).Int64()

		if basedAccountTotalAmount+feeAmount == TotalCetAmount {
			return nil
		}

		// Get all amounts charged by the commission
		return fmt.Errorf("The Cet total amount is inconsistent with the actual amount")
	}
}
