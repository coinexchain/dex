package bankx

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	dex "github.com/coinexchain/dex/types"
	"github.com/coinexchain/dex/x/authx"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgSetMemoRequired:
			return handleMsgSetMemoRequired(ctx, k.axk, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSend(ctx sdk.Context, k Keeper, msg bank.MsgSend) sdk.Result {

	if !k.bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.bk.Codespace()).Result()
	}

	fromAccount := k.ak.GetAccount(ctx, msg.FromAddress)
	amt := msg.Amount
	if !fromAccount.GetCoins().IsAllGTE(amt) {
		return sdk.ErrInsufficientCoins("sender has insufficient coins for the transfer").Result()
	}

	_, ok := k.axk.GetAccountX(ctx, msg.ToAddress)

	var neg bool
	activatedFee := k.GetParam(ctx).ActivatedFee

	//toaccount doesn't exist yet
	if !ok {
		//check whether the first transfer contains at least activatedFee cet
		amt, neg = amt.SafeSub(dex.NewCetCoins(activatedFee))
		if neg {
			return ErrorInsufficientCETForActivatingFee().Result()
		}

	}

	//handle coins transfer
	t, err := k.bk.SendCoins(ctx, msg.FromAddress, msg.ToAddress, amt)

	if err != nil {
		return err.Result()
	}

	// new accountx for toaddress if needed
	if !ok {
		// sub the activatedfees from fromaddress
		fromAccount = k.ak.GetAccount(ctx, msg.FromAddress)
		oldCoins := fromAccount.GetCoins()
		newCoins, _ := oldCoins.SafeSub(dex.NewCetCoins(activatedFee))
		fromAccount.SetCoins(newCoins)
		k.ak.SetAccount(ctx, fromAccount)

		//collect account activation fees
		k.fck.AddCollectedFees(ctx, dex.NewCetCoins(activatedFee))

		//create AccountX for toaddr
		newAccountX := authx.NewAccountXWithAddress(msg.ToAddress)
		newAccountX.Activated = true
		k.axk.SetAccountX(ctx, newAccountX)

	}

	return sdk.Result{
		Tags: t,
	}
}

func handleMsgSetMemoRequired(ctx sdk.Context, axk authx.AccountXKeeper, msg MsgSetMemoRequired) sdk.Result {
	accountX, found := axk.GetAccountX(ctx, msg.Address)
	if !found || !accountX.Activated {
		msg := fmt.Sprintf("account %s is not activated", msg.Address)
		return dex.ErrUnactivatedAddress(msg).Result()
	}

	accountX.TransferMemoRequired = msg.Required
	axk.SetAccountX(ctx, accountX)

	return sdk.Result{}
}
