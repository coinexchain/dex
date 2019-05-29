package bankx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/authx"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgSetMemoRequired:
			return handleMsgSetMemoRequired(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	if !k.bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.bk.Codespace()).Result()
	}

	fromAccount := k.ak.GetAccount(ctx, msg.FromAddress)
	toAccount := k.ak.GetAccount(ctx, msg.ToAddress)

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
		return sendLockedCoins(ctx, k, msg.FromAddress, msg.ToAddress, amt, time)
	}
	return normalSend(ctx, k, msg.FromAddress, msg.ToAddress, amt)
}

func deductActivationFee(ctx sdk.Context, k Keeper,
	fromAccount auth.Account, sendAmt sdk.Coins) (sdk.Coins, sdk.Error) {

	activatedFee := k.GetParam(ctx).ActivatedFee

	//check whether the first transfer contains at least activatedFee cet
	sendAmt, neg := sendAmt.SafeSub(dex.NewCetCoins(activatedFee))
	if neg {
		return sendAmt, ErrorInsufficientCETForActivatingFee()
	}

	// sub the activatedFees from fromAddress
	oldCoins := fromAccount.GetCoins()
	newCoins, _ := oldCoins.SafeSub(dex.NewCetCoins(activatedFee))
	fromAccount.SetCoins(newCoins)
	k.ak.SetAccount(ctx, fromAccount)

	//collect account activation fees
	k.fck.AddCollectedFees(ctx, dex.NewCetCoins(activatedFee))

	return sendAmt, nil
}

func sendLockedCoins(ctx sdk.Context, k Keeper,
	fromAddr, toAddr sdk.AccAddress, amt sdk.Coins, unlockTime int64) sdk.Result {

	ax := k.axk.GetOrCreateAccountX(ctx, toAddr)
	for _, coin := range amt {
		ax.LockedCoins = append(ax.LockedCoins, authx.LockedCoin{Coin: coin, UnlockTime: unlockTime})
	}
	k.axk.SetAccountX(ctx, ax)

	_, tag, err := k.bk.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tag,
	}
}

func normalSend(ctx sdk.Context, k Keeper,
	fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Result {

	t, err := k.bk.SendCoins(ctx, fromAddr, toAddr, amt)

	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: t,
	}
}

func handleMsgSetMemoRequired(ctx sdk.Context, k Keeper, msg MsgSetMemoRequired) sdk.Result {
	addr := msg.Address
	required := msg.Required

	account := k.ak.GetAccount(ctx, addr)
	if account == nil {
		msg := fmt.Sprintf("account %s is not activated", addr)
		return ErrUnactivatedAddress(msg).Result()
	}

	accountX := k.axk.GetOrCreateAccountX(ctx, addr)
	accountX.MemoRequired = required
	k.axk.SetAccountX(ctx, accountX)

	return sdk.Result{
		Tags: sdk.NewTags(
			TagKeyMemoRequired, fmt.Sprintf("%v", required),
			TagKeyAddr, addr.String(),
		),
	}
}
