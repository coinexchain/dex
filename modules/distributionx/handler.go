package distributionx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"

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
	res := k.bxk.Sk.SendCoinsFromAccountToModule(ctx, msg.FromAddr, distribution.ModuleName, msg.Amount)
	if res != nil {
		return res.Result()
	}

	k.AddCoinsToFeePool(ctx, msg.Amount)

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	TagKeyDonator, msg.FromAddr.String(),
		//),
	}
}
