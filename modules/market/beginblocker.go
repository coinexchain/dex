package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

func BeginBlocker(ctx sdk.Context, k keepers.Keeper) {
	msg := types.NewHeightInfo{
		Height:    ctx.BlockHeight(),
		TimeStamp: ctx.BlockHeader().Time.Unix(),
	}

	k.SendMsg(types.HeightInfoKey, msg)
}
