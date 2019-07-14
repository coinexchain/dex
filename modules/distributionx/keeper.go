package distributionx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

type Keeper struct {
	bxk bankx.Keeper
	dk  distribution.Keeper
}

func NewKeeper(bxk bankx.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		bxk,
		dk,
	}
}

func (keeper Keeper) AddCoinsToFeePool(ctx sdk.Context, coins sdk.Coins) {

	feePool := keeper.dk.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(coins))
	keeper.dk.SetFeePool(ctx, feePool)

}
