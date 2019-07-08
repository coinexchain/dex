package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	msg := NewHeightInfo{
		Height:    ctx.BlockHeight(),
		TimeStamp: ctx.BlockHeader().Time.Unix(),
	}
	k.SendMsg(HeightInfoKey, msg)
}
