package bankx

import (
	"fmt"
	types2 "github.com/coinexchain/dex/modules/authx/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/bankx/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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

	err := k.Sk.SendCoinsFromAccountToModule(ctx, fromAddr, types2.ModuleName, amt)
	if err != nil {
		return err.Result()
	}

	ax := k.Axk.GetOrCreateAccountX(ctx, toAddr)
	for _, coin := range amt {
		ax.LockedCoins = append(ax.LockedCoins, types2.LockedCoin{Coin: coin, UnlockTime: unlockTime})
	}
	k.Axk.SetAccountX(ctx, ax)

	if amt != nil {
		k.Axk.InsertUnlockedCoinsQueue(ctx, unlockTime, toAddr)
	}

	e := k.MsgProducer.SendMsg(types.Topic, "send_lock_coins", types.NewMsgSend(fromAddr, toAddr, amt, unlockTime))
	if e != nil {
		//record fail sendMsg log
	}
	return sdk.Result{
		//Tags: tag,
	}
}

func normalSend(ctx sdk.Context, k Keeper,
	fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Result {

	err := k.Bk.SendCoins(ctx, fromAddr, toAddr, amt)

	if err != nil {
		return err.Result()
	}

	e := k.MsgProducer.SendMsg(types.Topic, "send_coins", types.NewMsgSend(fromAddr, toAddr, amt, 0))
	if e != nil {
		//record fail sendMsg log
	}

	return sdk.Result{
		//Tags: t,
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
	e := k.MsgProducer.SendMsg(types.Topic, "send_memo_required", msg)
	if e != nil {
		//record fail sendMsg log
	}
	return sdk.Result{
		//Tags: sdk.NewTags(
		//	TagKeyMemoRequired, fmt.Sprintf("%v", required),
		//	TagKeyAddr, addr.String(),
		//),
	}
}
