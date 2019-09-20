package comment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/comment/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCommentToken:
			return handleMsgCommentToken(ctx, k, msg)
		default:
			errMsg := "Unrecognized comment Msg type: " + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCommentToken(ctx sdk.Context, k Keeper, msg types.MsgCommentToken) sdk.Result {
	if !k.Axk.IsTokenExists(ctx, msg.Token) {
		return types.ErrNoSuchAsset().Result()
	}
	if msg.Donation > 0 {
		donatedCoin := sdk.Coins{sdk.NewCoin("cet", sdk.NewInt(msg.Donation))}
		res := k.Dk.DonateToCommunityPool(ctx, msg.Sender, donatedCoin)
		if res != nil {
			return res.Result()
		}
	}

	for _, ref := range msg.References {
		if ref.RewardAmount <= 0 {
			continue
		}
		rewardCoin := sdk.NewCoin(ref.RewardToken, sdk.NewInt(ref.RewardAmount))
		acc := k.Ak.GetAccount(ctx, ref.RewardTarget)
		if acc == nil {
			return types.ErrNoSuchAccount(ref.RewardTarget.String()).Result()
		}
		res := k.Bxk.SendCoins(ctx, msg.Sender, ref.RewardTarget, sdk.Coins{rewardCoin})
		if res != nil {
			return res.Result()
		}
	}

	lastCount := k.Cck.IncrCommentCount(ctx, msg.Token)

	if len(k.EventTypeMsgQueue) != 0 {
		tokenComment := types.NewTokenComment(&msg, lastCount, ctx.BlockHeight())
		bytes := dex.SafeJSONMarshal(tokenComment)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				k.EventTypeMsgQueue,
				sdk.NewAttribute(types.TokenCommentKey, string(bytes)),
			),
		)
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
