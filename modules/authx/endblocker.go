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
			transferUnlockedCoins(&acc, currentTime, ctx, aux, keeper)
			aux.RemoveFromUnlockedCoinsQueueByKey(ctx, iterator.Key())
		}
	}
}

func transferUnlockedCoins(accx *types.AccountX, time int64, ctx sdk.Context, kx AccountXKeeper, keeper ExpectedAccountKeeper) {
	var coins = sdk.Coins{}
	var temp types.LockedCoins
	for _, c := range accx.LockedCoins {
		if c.UnlockTime <= time {
			coins = coins.Add(sdk.Coins{c.Coin})
		} else {
			temp = append(temp, c)
		}
	}
	coins = coins.Sort()

	acc := keeper.GetAccount(ctx, accx.Address)
	newCoins := acc.GetCoins().Add(coins)
	_ = acc.SetCoins(newCoins)
	keeper.SetAccount(ctx, acc)

	accx.LockedCoins = temp
	kx.SetAccountX(ctx, *accx)
}
