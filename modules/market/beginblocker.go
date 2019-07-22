package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

func BeginBlocker(ctx sdk.Context, k keepers.Keeper) {
	if k.IsSubScribe(types.Topic) {
		msg := types.NewHeightInfo{
			Height:    ctx.BlockHeight(),
			TimeStamp: ctx.BlockHeader().Time.Unix(),
		}
		fillMsgs(ctx, types.HeightInfoKey, msg)
	}
}
