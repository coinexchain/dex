package comment

import (
	"github.com/coinexchain/dex/modules/comment/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCommentToken:
			return handleMsgCommentToken(ctx, k, msg)
		default:
			errMsg := "Unrecognized comment Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCommentToken(ctx sdk.Context, k Keeper, msg types.MsgCommentToken) sdk.Result {
	if !k.Axk.IsTokenExists(ctx, msg.Token) {
		return types.ErrNoSuchAsset().Result()
	}
	if msg.Donation > 0 {
		donatedCoin := sdk.Coins{sdk.Coin{Denom: "cet", Amount: sdk.NewInt(msg.Donation)}}
		res := k.Dk.DonateToCommunityPool(ctx, msg.Sender, donatedCoin)
		if res != nil {
			return res.Result()
		}
	}

	for _, ref := range msg.References {
		if ref.RewardAmount <= 0 {
			continue
		}
		rewardCoin := sdk.Coin{Denom: ref.RewardToken, Amount: sdk.NewInt(ref.RewardAmount)}
		res := k.Bxk.SendCoins(ctx, msg.Sender, ref.RewardTarget, sdk.Coins{rewardCoin})
		if res != nil {
			return res.Result()
		}
	}

	lastCount := k.Cck.IncrCommentCount(ctx, msg.Token)

	if k.MsgSendFunc != nil {
		tokenComment := types.NewTokenComment(&msg, lastCount)
		k.MsgSendFunc(types.TokenCommentKey, tokenComment)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return sdk.Result{
		Codespace: types.CodeSpaceComment,
		Events:    ctx.EventManager().Events(),
	}
}
