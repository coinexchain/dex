package crisisx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/coinexchain/dex/types"
)

const (
	TotalCetAmount int64 = 588800000000000000
	ModuleName           = "crisisx"
)

func RegisterInvariants(c *crisis.Keeper, bk ExpectBankxKeeper, feek auth.FeeCollectionKeeper, disk distribution.Keeper) {
	c.RegisterRoute(ModuleName, "cet-invariant", SupplyCETInvariant(bk, feek, disk))
}

func SupplyCETInvariant(bk ExpectBankxKeeper, feek auth.FeeCollectionKeeper, disk distribution.Keeper) sdk.Invariant {

	return func(ctx sdk.Context) error {

		// Get all amounts based on the account system
		basedAccountTotalAmount := bk.TotalAmountOfCoin(ctx, types.CET)
		feeAmount := feek.GetCollectedFees(ctx).AmountOf(types.CET).Int64()
		communityAmount := disk.GetFeePool(ctx).CommunityPool.AmountOf(types.CET).Int64()
		fmt.Printf("basedAccountTotalAmount : %d, feeAmount : %d, communityAmount : %d\n",
			basedAccountTotalAmount, feeAmount, communityAmount)
		if basedAccountTotalAmount+feeAmount+communityAmount == TotalCetAmount {
			return nil
		}

		return fmt.Errorf("The Cet total amount is inconsistent with the actual amount")
	}
}
