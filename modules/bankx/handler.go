package bankx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)
		case types.MsgSetMemoRequired:
			return handleMsgSetMemoRequired(ctx, k, msg)
		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)
		case types.MsgSupervisedSend:
			return handleMsgSupervisedSend(ctx, k, msg)
		default:
			return dex.ErrUnknownRequest(ModuleName, msg)
		}
	}
}

func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg types.MsgMultiSend) sdk.Result {
	if enabled := k.GetSendEnabled(ctx); !enabled {
		return bank.ErrSendDisabled(types.CodeSpaceBankx).Result()
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}
	}

	for _, input := range msg.Inputs {
		if k.IsSendForbidden(ctx, input.Coins, input.Address) {
			return types.ErrTokenForbiddenByOwner().Result()
		}
	}

	addrs := k.PreCheckFreshAccounts(ctx, msg.Outputs)

	if err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs); err != nil {
		return err.Result()
	}

	if err := k.DeductActivationFeeForFreshAccounts(ctx, addrs); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	//fillMsgQueue(ctx, k, "multi_send", msg)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}

}

func handleMsgSend(ctx sdk.Context, k Keeper, msg types.MsgSend) sdk.Result {
	if enabled := k.GetSendEnabled(ctx); !enabled {
		return bank.ErrSendDisabled(types.CodeSpaceBankx).Result()
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
	}

	if k.IsSendForbidden(ctx, msg.Amount, msg.FromAddress) {
		return types.ErrTokenForbiddenByOwner().Result()
	}

	//TODO: add codes to check whether fromAccount & toAccount is moduleAccount

	amt := msg.Amount
	if !k.GetCoins(ctx, msg.FromAddress).IsAllGTE(amt) {
		return sdk.ErrInsufficientCoins("sender has insufficient coins for the transfer").Result()
	}

	//check whether toAccount exist
	amt, err := k.DeductActivationFee(ctx, msg.FromAddress, msg.ToAddress, amt)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	time := msg.UnlockTime
	if time != 0 {
		return lockedSend(ctx, k, msg.FromAddress, msg.ToAddress, amt, time)
	}
	return normalSend(ctx, k, msg.FromAddress, msg.ToAddress, amt)
}

func lockedSend(ctx sdk.Context, k Keeper, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins, unlockTime int64) sdk.Result {
	if unlockTime < ctx.BlockHeader().Time.Unix() {
		return types.ErrUnlockTime("Invalid Unlock Time:" +
			fmt.Sprintf("%d < %d", unlockTime, ctx.BlockHeader().Time.Unix())).Result()
	}

	if err := k.SendLockedCoins(ctx, fromAddr, toAddr, nil, amt, unlockTime, 0, false); err != nil {
		return err.Result()
	}

	fillMsgQueue(ctx, k, "send_lock_coins", types.NewLockedSendMsg(fromAddr, toAddr, amt, unlockTime))

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func normalSend(ctx sdk.Context, k Keeper, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Result {
	err := k.SendCoins(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err.Result()
	}

	//fillMsgQueue(ctx, k, "send_coins", types.NewMsgSend(fromAddr, toAddr, amt, 0))

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSetMemoRequired(ctx sdk.Context, k Keeper, msg types.MsgSetMemoRequired) sdk.Result {
	addr := msg.Address
	required := msg.Required

	if k.BlacklistedAddr(msg.Address) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.Address)).Result()
	}

	if err := k.SetMemoRequired(ctx, addr, required); err != nil {
		return err.Result()
	}

	//fillMsgQueue(ctx, k, "send_memo_required", msg)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAddr, addr.String()),
			sdk.NewAttribute(types.AttributeKeyMemoRequired, fmt.Sprintf("%v", required)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSupervisedSend(ctx sdk.Context, k Keeper, msg types.MsgSupervisedSend) sdk.Result {
	if enabled := k.GetSendEnabled(ctx); !enabled {
		return bank.ErrSendDisabled(types.CodeSpaceBankx).Result()
	}

	if k.GetAccount(ctx, msg.ToAddress) == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", msg.ToAddress)).Result()
	}
	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
	}
	if !msg.Supervisor.Empty() && msg.Reward > 0 {
		if k.GetAccount(ctx, msg.Supervisor) == nil {
			return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", msg.Supervisor)).Result()
		}
		if k.BlacklistedAddr(msg.Supervisor) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.Supervisor)).Result()
		}
	}

	if msg.UnlockTime < ctx.BlockHeader().Time.Unix() {
		return types.ErrUnlockTime("Invalid Unlock Time:" + fmt.Sprintf("%d < %d", msg.UnlockTime, ctx.BlockHeader().Time.Unix())).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	if msg.Operation == types.Create {
		amt := sdk.NewCoins(msg.Amount)
		if !k.GetCoins(ctx, msg.FromAddress).IsAllGTE(amt) {
			return sdk.ErrInsufficientCoins("sender has insufficient coin for the transfer").Result()
		}
		if err := k.SendLockedCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Supervisor, amt, msg.UnlockTime, msg.Reward, true); err != nil {
			return err.Result()
		}

		fillMsgQueue(ctx, k, "send_lock_coins", types.NewSupervisedSendMsg(msg.FromAddress, msg.ToAddress,
			msg.Supervisor, msg.Amount, msg.UnlockTime, msg.Reward))

	} else {
		isReturned := types.Return == msg.Operation
		if err := k.EarlierUnlockCoin(ctx, msg.FromAddress, msg.ToAddress, msg.Supervisor, &msg.Amount, msg.UnlockTime, msg.Reward, isReturned); err != nil {
			return err.Result()
		}

		acc := k.GetAccount(ctx, msg.ToAddress)
		accx, _ := k.GetAccountX(ctx, msg.ToAddress)
		fillMsgQueue(ctx, k, "notify_unlock", authx.NotificationUnlock{
			Address:     msg.ToAddress,
			Unlocked:    sdk.NewCoins(msg.Amount),
			LockedCoins: accx.LockedCoins,
			FrozenCoins: accx.FrozenCoins,
			Coins:       acc.GetCoins(),
			Height:      ctx.BlockHeight(),
		})
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
	if keeper.MsgProducer.IsSubscribed(types.Topic) {
		msgqueue.FillMsgs(ctx, key, msg)
	}
}
