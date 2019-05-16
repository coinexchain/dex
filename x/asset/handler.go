package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// handleMsgIssueToken - Handle MsgIssueToken
func handleMsgIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
	_, err := k.IssueToken(ctx, msg)
	if err != nil {
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
