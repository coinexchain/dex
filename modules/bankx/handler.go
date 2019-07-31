package bankx

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
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
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSend(ctx sdk.Context, k Keeper, msg types.MsgSend) sdk.Result {
	if !k.Bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Bk.Codespace()).Result()
	}

	if k.IsSendForbidden(ctx, msg.Amount, msg.FromAddress) {
		return types.ErrTokenForbiddenByOwner("transfer has been forbidden by token owner").Result()
	}

	fromAccount := k.Ak.GetAccount(ctx, msg.FromAddress)
	toAccount := k.Ak.GetAccount(ctx, msg.ToAddress)

	//TODO: add codes to check whether fromAccount & toAccount is moduleAccount

	amt := msg.Amount
	if !fromAccount.GetCoins().IsAllGTE(amt) {
		return sdk.ErrInsufficientCoins("sender has insufficient coins for the transfer").Result()
	}

	//toAccount doesn't exist yet
	if toAccount == nil {
		var err sdk.Error
		amt, err = deductActivationFee(ctx, k, fromAccount, amt)
		if err != nil {
			return err.Result()
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	time := msg.UnlockTime
	if time != 0 {
		if time < ctx.BlockHeader().Time.Unix() {
			return types.ErrUnlockTime("Invalid Unlock Time").Result()
		}
		return sendLockedCoins(ctx, k, msg.FromAddress, msg.ToAddress, amt, time)
	}
	return normalSend(ctx, k, msg.FromAddress, msg.ToAddress, amt)
}

func deductActivationFee(ctx sdk.Context, k Keeper,
	fromAccount auth.Account, sendAmt sdk.Coins) (sdk.Coins, sdk.Error) {

	activationFee := k.GetParam(ctx).ActivationFee

	//check whether the first transfer contains at least activationFee cet
	sendAmt, neg := sendAmt.SafeSub(dex.NewCetCoins(activationFee))
	if neg {
		return sendAmt, types.ErrorInsufficientCETForActivatingFee()
	}

	err := k.DeductFee(ctx, fromAccount.GetAddress(), dex.NewCetCoins(activationFee))
	if err != nil {
		return sendAmt, err
	}

	return sendAmt, nil
}

func sendLockedCoins(ctx sdk.Context, k Keeper,
	fromAddr, toAddr sdk.AccAddress, amt sdk.Coins, unlockTime int64) sdk.Result {

	if k.Ak.GetAccount(ctx, toAddr) == nil {
		_, _ = k.Bk.AddCoins(ctx, toAddr, sdk.Coins{})
	}

	if err := k.DeductFee(ctx, fromAddr, dex.NewCetCoins(k.GetParam(ctx).LockCoinsFee)); err != nil {
		return err.Result()
	}

	_, err := k.Bk.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err.Result()
	}

	ax := k.Axk.GetOrCreateAccountX(ctx, toAddr)
	for _, coin := range amt {
		ax.LockedCoins = append(ax.LockedCoins, authx.LockedCoin{Coin: coin, UnlockTime: unlockTime})
	}
	k.Axk.SetAccountX(ctx, ax)

	if amt != nil {
		k.Axk.InsertUnlockedCoinsQueue(ctx, unlockTime, toAddr)
	}

	fillMsgQueue(ctx, k, "send_lock_coins", types.NewMsgSend(fromAddr, toAddr, amt, unlockTime))

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func normalSend(ctx sdk.Context, k Keeper,
	fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Result {

	err := k.Bk.SendCoins(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err.Result()
	}
	fillMsgQueue(ctx, k, "send_coins", types.NewMsgSend(fromAddr, toAddr, amt, 0))

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSetMemoRequired(ctx sdk.Context, k Keeper, msg types.MsgSetMemoRequired) sdk.Result {
	addr := msg.Address
	required := msg.Required

	account := k.Ak.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr)).Result()
	}

	accountX := k.Axk.GetOrCreateAccountX(ctx, addr)
	accountX.MemoRequired = required
	k.Axk.SetAccountX(ctx, accountX)
	fillMsgQueue(ctx, k, "send_memo_required", msg)

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

func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
	if keeper.MsgProducer.IsSubScribe(types.Topic) {
		bytes, err := json.Marshal(msg)
		if err != nil {
			return
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(msgqueue.EventTypeMsgQueue,
			sdk.NewAttribute(key, string(bytes))))
	}
}
