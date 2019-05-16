package asset

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// NewHandler returns a handler for "asset" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssueToken:
			return handleMsgIssueToken(ctx, k, msg)
		case MsgTransferOwnership:
			return handleMsgTransferOwnership(ctx, k, msg)
		case MsgFreezeAddress:
			return handleMsgFreezeAddress(ctx, k, msg)
		case MsgUnfreezeAddress:
			return handleMsgUnfreezeAddress(ctx, k, msg)
		case MsgFreezeToken:
			return handleMsgFreezeToken(ctx, k, msg)
		case MsgUnfreezeToken:
			return handleMsgUnfreezeToken(ctx, k, msg)
		case MsgBurnToken:
			return handleMsgBurnToken(ctx, k, msg)
		case MsgMintToken:
			return handleMsgMintToken(ctx, k, msg)

		default:
			errMsg := "Unrecognized asset Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress("no issue address")
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

func subTokenFee(ctx sdk.Context, k Keeper, addr sdk.AccAddress, fee sdk.Coins) sdk.Error {

	acc := k.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress("no valid address")
	}

	oldCoins := acc.GetCoins()
	spendableCoins := acc.SpendableCoins(ctx.BlockHeader().Time)

	_, hasNeg := spendableCoins.SafeSub(fee)
	if hasNeg {
		return sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, fee))
	}

	newCoins := oldCoins.Sub(fee) // should not panic as spendable coins was already checked
	if err := setCoins(ctx, k.ak, addr, newCoins); err != nil {
		return err
	}

	return nil
}


// handleMsgIssueToken - Handle MsgIssueToken
func handleMsgIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {

	issueFee := k.GetParams(ctx).IssueTokenFee
	if err := subTokenFee(ctx, k, msg.Owner, issueFee); err != nil {
		return err.Result()
	}
	k.fck.AddCollectedFees(ctx, issueFee)

	if err := k.IssueToken(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handleMsgTransferOwnership - Handle MsgTransferOwnership
func handleMsgTransferOwnership(ctx sdk.Context, k Keeper, msg MsgTransferOwnership) (res sdk.Result) {

	return
}

// handleMsgFreezeAddress - Handle MsgFreezeAddress
func handleMsgFreezeAddress(ctx sdk.Context, k Keeper, msg MsgFreezeAddress) (res sdk.Result) {

	return
}

// handleMsgUnfreezeAddress - Handle MsgUnfreezeAddress
func handleMsgUnfreezeAddress(ctx sdk.Context, k Keeper, msg MsgUnfreezeAddress) (res sdk.Result) {

	return
}

// handleMsgFreezeToken - HandleMsgFreezeToken
func handleMsgFreezeToken(ctx sdk.Context, k Keeper, msg MsgFreezeToken) (res sdk.Result) {

	return
} // handleMsgUnfreezeToken - Handle MsgUnfreezeToken
func handleMsgUnfreezeToken(ctx sdk.Context, k Keeper, msg MsgUnfreezeToken) (res sdk.Result) {

	return
}

// handleMsgBurnToken - Handle MsgBurnToken
func handleMsgBurnToken(ctx sdk.Context, k Keeper, msg MsgBurnToken) (res sdk.Result) {

	return
}

// handleMsgMintToken - Handle MsgMintToken
func handleMsgMintToken(ctx sdk.Context, k Keeper, msg MsgMintToken) (res sdk.Result) {

	return
}
