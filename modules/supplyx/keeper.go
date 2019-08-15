package supplyx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// "extends" supply keeper to "override" BurnCoins() method
type Keeper struct {
	supply.Keeper
	dk distribution.Keeper
}

func NewKeeper(sk supply.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		Keeper: sk,
		dk:     dk,
	}
}

func (k Keeper) BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error {
	err := k.Keeper.SendCoinsFromModuleToModule(ctx, name, distribution.ModuleName, amt)
	if err != nil {
		return err
	}

	feePool := k.dk.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(amt))
	k.dk.SetFeePool(ctx, feePool)
	return nil
}
