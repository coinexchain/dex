package distributionx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/distributionx/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgDonateToCommunityPool:
			return handleMsgDonateToCommunityPool(ctx, k, msg)
		default:
			errMsg := "Unrecognized distributionx Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDonateToCommunityPool(ctx sdk.Context, k Keeper, msg types.MsgDonateToCommunityPool) sdk.Result {

	err := k.DonateToCommunityPool(ctx, msg.FromAddr, msg.Amount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		//Tags: sdk.NewTags(
		//	TagKeyDonator, msg.FromAddr.String(),
		//),
	}
}
