package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/authx/types"
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
			withdrawUnlockedCoins(&acc, currentTime, ctx, aux, keeper)
			aux.RemoveFromUnlockedCoinsQueueByKey(ctx, iterator.Key())
		}
	}
}

func withdrawUnlockedCoins(accx *types.AccountX, time int64, ctx sdk.Context, kx AccountXKeeper, keeper ExpectedAccountKeeper) {
	var unlocked = sdk.Coins{}
	var stillLocked types.LockedCoins
	for _, c := range accx.LockedCoins {
		if c.UnlockTime <= time {
			unlocked = unlocked.Add(sdk.Coins{c.Coin})
		} else {
			stillLocked = append(stillLocked, c)
		}
	}
	unlocked = unlocked.Sort()

	acc := keeper.GetAccount(ctx, accx.Address)
	newCoins := acc.GetCoins().Add(unlocked)
	_ = acc.SetCoins(newCoins)
	keeper.SetAccount(ctx, acc)

	accx.LockedCoins = stillLocked
	kx.SetAccountX(ctx, *accx)
}
