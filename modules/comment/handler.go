package comment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCommentToken:
			return handleMsgCommentToken(ctx, k, msg)
		default:
			errMsg := "Unrecognized comment Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCommentToken(ctx sdk.Context, k Keeper, msg MsgCommentToken) sdk.Result {
	if !k.axk.IsTokenExists(ctx, msg.Token) {
		return ErrNoSuchAsset().Result()
	}
	if msg.Donation > 0 {
		donatedCoin := sdk.Coins{sdk.Coin{Denom: "cet", Amount: sdk.NewInt(msg.Donation)}}
		k.dk.DonateToCommunityPool(ctx, msg.Sender, donatedCoin)
	}

	for _, ref := range msg.References {
		if ref.RewardAmount <= 0 {
			continue
		}
		rewardCoin := sdk.Coin{Denom: ref.RewardToken, Amount: sdk.NewInt(ref.RewardAmount)}
		res := k.bxk.SendCoins(ctx, msg.Sender, ref.RewardTarget, sdk.Coins{rewardCoin})
		if res != nil {
			return res.Result()
		}
	}

	if k.msgSendFunc != nil {
		tokenComment := NewTokenComment(&msg, k.cck.GetCommentCount(ctx))
		k.msgSendFunc(TokenCommentKey, tokenComment)
	}

	k.cck.IncrCommentCount(ctx)

	return sdk.Result{
		Codespace: CodeSpaceComment,
		//Tags: sdk.NewTags(
		//	distributionx.TagKeyDonator, msg.Sender.String(),
		//),
	}
}
