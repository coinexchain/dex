package stakingx

import (
	"fmt"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TotalSupplyInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		token := k.assetViewKeeper.GetToken(ctx, types.DefaultBondDenom)
		if token == nil {
			return fmt.Errorf("cet not found")
		}

		//ts := token.GetTotalSupply()
		//bondPool := k.CalcBondPoolStatus(ctx)

		//TODO: compare with supplyKeeper.TotalSupply
		//if ts != bondPool.TotalSupply.Int64() {
		//	return fmt.Errorf("total-supply invariance:\n"+
		//		"\tinconsistent total-supply: \n"+
		//		"\tCET asset total supply: %v\n"+
		//		"\tpool: %v", ts, bondPool)
		//}

		return nil
	}
}

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
