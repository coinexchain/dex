package stakingx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	dex "github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type AssetViewKeeper interface {
	GetToken(ctx sdk.Context, symbol string) asset.Token
}

func RegisterInvariants(c types.CrisisKeeper, k Keeper, assetKeeper AssetViewKeeper, sk staking.Keeper) {

	c.RegisterRoute(types.ModuleName, "total-supply", TotalSupplyInvariants(k, assetKeeper))

	// SupplyInvariants no longer suitable here, new SupplyInvariants will be created
	// c.RegisterRoute(types.ModuleName, "supply",
	//	SupplyInvariants(k, f, d, am))

	c.RegisterRoute(types.ModuleName, "nonnegative-power",
		staking.NonNegativePowerInvariant(sk))

	c.RegisterRoute(types.ModuleName, "positive-delegation",
		staking.PositiveDelegationInvariant(sk))

	c.RegisterRoute(types.ModuleName, "delegator-shares",
		staking.DelegatorSharesInvariant(sk))
}

func TotalSupplyInvariants(k Keeper, assetKeeper AssetViewKeeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		token := assetKeeper.GetToken(ctx, dex.DefaultBondDenom)
		if token == nil {
			return nil
		}

		ts := token.GetTotalSupply()
		bondPool := k.CalcBondPoolStatus(ctx)

		if ts != bondPool.TotalSupply.Int64() {
			return fmt.Errorf("total-supply invariance:\n"+
				"\tinconsistent total-supply: \n"+
				"\tCET asset total supply: %v\n"+
				"\tpool: %v", ts, bondPool)
		}

		return nil
	}
}
