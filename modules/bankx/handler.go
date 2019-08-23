package bankx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
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
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg types.MsgMultiSend) sdk.Result {
	if !k.Bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Bk.Codespace()).Result()
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}
	}

	if err := checkInputs(ctx, k, msg.Inputs); err != nil {
		return err.Result()
	}

	addrs := precheckFreshAccounts(ctx, k, msg.Outputs)

	if err := k.Bk.InputOutputCoins(ctx, msg.Inputs, msg.Outputs); err != nil {
		return err.Result()
	}

	if err := deductActivationFeeForFreshAccount(ctx, k, addrs); err != nil {
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
	if !k.Bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Bk.Codespace()).Result()
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
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
			return types.ErrUnlockTime("Invalid Unlock Time:" + fmt.Sprintf("%d < %d", time, ctx.BlockHeader().Time.Unix())).Result()
		}
		return sendLockedCoins(ctx, k, msg.FromAddress, msg.ToAddress, amt, time)
	}
	return normalSend(ctx, k, msg.FromAddress, msg.ToAddress, amt)
}

func deductActivationFee(ctx sdk.Context, k Keeper,
	fromAccount auth.Account, sendAmt sdk.Coins) (sdk.Coins, sdk.Error) {

	activationFee := k.GetParams(ctx).ActivationFee

	//check whether the first transfer contains at least activationFee cet
	sendAmt, neg := sendAmt.SafeSub(dex.NewCetCoins(activationFee))
	if neg {
		return sendAmt, types.ErrorInsufficientCETForActivatingFee()
	}

	err := k.DeductInt64CetFee(ctx, fromAccount.GetAddress(), activationFee)
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

	if err := k.DeductInt64CetFee(ctx, fromAddr, k.GetParams(ctx).LockCoinsFee); err != nil {
		return err.Result()
	}

	_, err := k.Bk.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err.Result()
	}

	ax := k.Axk.GetOrCreateAccountX(ctx, toAddr)
	for _, coin := range amt {
		ax.LockedCoins = append(ax.LockedCoins, authx.LockedCoin{Coin: coin, UnlockTime: unlockTime})
		if err := k.Tk.UpdateTokenSendLock(ctx, coin.Denom, coin.Amount, true); err != nil {
			return err.Result()
		}
	}
	k.Axk.SetAccountX(ctx, ax)

	if !amt.Empty() {
		k.Axk.InsertUnlockedCoinsQueue(ctx, unlockTime, toAddr)
	}

	//fillMsgQueue(ctx, k, "send_lock_coins", types.NewMsgSend(fromAddr, toAddr, amt, unlockTime))

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

	account := k.Ak.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr)).Result()
	}

	accountX := k.Axk.GetOrCreateAccountX(ctx, addr)
	accountX.MemoRequired = required
	k.Axk.SetAccountX(ctx, accountX)
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

//func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
//	if keeper.MsgProducer.IsSubscribed(types.Topic) {
//		bytes, err := json.Marshal(msg)
//		if err != nil {
//			return
//		}
//		ctx.EventManager().EmitEvent(sdk.NewEvent(msgqueue.EventTypeMsgQueue,
//			sdk.NewAttribute(key, string(bytes))))
//	}
//}

func checkInputs(ctx sdk.Context, k Keeper, inputs []bank.Input) sdk.Error {
	for _, input := range inputs {
		if k.IsSendForbidden(ctx, input.Coins, input.Address) {
			return types.ErrTokenForbiddenByOwner("transfer has been forbidden by token owner")
		}

		fromAccount := k.Ak.GetAccount(ctx, input.Address)
		if !fromAccount.GetCoins().IsAllGTE(input.Coins) {
			return sdk.ErrInsufficientCoins("sender has insufficient coins for the transfer")
		}

	}
	return nil

}

func precheckFreshAccounts(ctx sdk.Context, k Keeper, outputs []bank.Output) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	for _, output := range outputs {
		toAccount := k.Ak.GetAccount(ctx, output.Address)

		//toAccount doesn't exist yet
		if toAccount == nil {
			addrs = append(addrs, output.Address)
		}
	}
	return addrs
}
func deductActivationFeeForFreshAccount(ctx sdk.Context, k Keeper, addrs []sdk.AccAddress) sdk.Error {
	for _, addr := range addrs {
		activationFee := k.GetParams(ctx).ActivationFee
		err := k.DeductInt64CetFee(ctx, addr, activationFee)
		if err != nil {
			return err
		}
	}

	return nil
}
