package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, aux AccountXKeeper, keeper ExpectedAccountKeeper) {
	currentTime := ctx.BlockHeader().Time.Unix()
	iterator := aux.UnlockedCoinsQueueIterator(ctx, currentTime)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := iterator.Value()
		if addr != nil {
			acc, ok := aux.GetAccountX(ctx, addr)
			if !ok {
				//always account exist
				continue
			}
			acc.TransferUnlockedCoins(currentTime, ctx, aux, keeper)
			aux.RemoveFromUnlockedCoinsQueueByKey(ctx, iterator.Key())
		}
	}
}
