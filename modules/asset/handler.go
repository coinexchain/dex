package asset

import (
	"strconv"

	"github.com/coinexchain/dex/modules/asset/tags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "asset" type messages.
func NewHandler(tk TokenKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgIssueToken:
			return handleMsgIssueToken(ctx, tk, msg)
		case MsgTransferOwnership:
			return handleMsgTransferOwnership(ctx, tk, msg)
		case MsgMintToken:
			return handleMsgMintToken(ctx, tk, msg)
		case MsgBurnToken:
			return handleMsgBurnToken(ctx, tk, msg)
		case MsgForbidToken:
			return handleMsgForbidToken(ctx, tk, msg)
		case MsgUnForbidToken:
			return handleMsgUnForbidToken(ctx, tk, msg)
		case MsgAddTokenWhitelist:
			return handleMsgAddTokenWhitelist(ctx, tk, msg)
		case MsgRemoveTokenWhitelist:
			return handleMsgRemoveTokenWhitelist(ctx, tk, msg)
		case MsgForbidAddr:
			return handleMsgForbidAddr(ctx, tk, msg)
		case MsgUnForbidAddr:
			return handleMsgUnForbidAddr(ctx, tk, msg)

		default:
			errMsg := "Unrecognized asset Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgIssueToken - Handle MsgIssueToken
func handleMsgIssueToken(ctx sdk.Context, tk TokenKeeper, msg MsgIssueToken) sdk.Result {

	issueFee := tk.GetParams(ctx).IssueTokenFee
	if err := tk.SubtractFeeAndCollectFee(ctx, msg.Owner, issueFee); err != nil {
		return err.Result()
	}

	if err := tk.IssueToken(ctx, msg); err != nil {
		return err.Result()
	}

	if err := tk.AddToken(ctx, msg.Owner, NewTokenCoins(msg.Symbol, msg.TotalSupply)); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.Owner, msg.Owner.String(),
		),
	}
}

// handleMsgTransferOwnership - Handle MsgTransferOwnership
func handleMsgTransferOwnership(ctx sdk.Context, tk TokenKeeper, msg MsgTransferOwnership) sdk.Result {
	if err := tk.TransferOwnership(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.OriginalOwner, msg.OriginalOwner.String(),
			tags.NewOwner, msg.NewOwner.String(),
		),
	}
}

// handleMsgMintToken - Handle MsgMintToken
func handleMsgMintToken(ctx sdk.Context, tk TokenKeeper, msg MsgMintToken) sdk.Result {
	if err := tk.MintToken(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.Amt, strconv.FormatInt(msg.Amount, 10),
		),
	}
}

// handleMsgBurnToken - Handle MsgBurnToken
func handleMsgBurnToken(ctx sdk.Context, tk TokenKeeper, msg MsgBurnToken) sdk.Result {
	if err := tk.BurnToken(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.Amt, strconv.FormatInt(msg.Amount, 10),
		),
	}
}

// handleMsgForbidToken - Handle ForbidToken msg
func handleMsgForbidToken(ctx sdk.Context, tk TokenKeeper, msg MsgForbidToken) sdk.Result {
	if err := tk.ForbidToken(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
		),
	}
}

// handleMsgUnForbidToken - Handle UnForbidToken msg
func handleMsgUnForbidToken(ctx sdk.Context, tk TokenKeeper, msg MsgUnForbidToken) sdk.Result {
	if err := tk.UnForbidToken(ctx, msg); err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
		),
	}
}

// handleMsgAddTokenWhitelist - Handle AddTokenWhitelist msg
func handleMsgAddTokenWhitelist(ctx sdk.Context, tk TokenKeeper, msg MsgAddTokenWhitelist) sdk.Result {
	if err := tk.AddTokenWhitelist(ctx, msg); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String() + ","
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.AddWhitelist, str,
		),
	}
}

// handleMsgRemoveTokenWhitelist - Handle RemoveTokenWhitelist msg
func handleMsgRemoveTokenWhitelist(ctx sdk.Context, tk TokenKeeper, msg MsgRemoveTokenWhitelist) sdk.Result {
	if err := tk.RemoveTokenWhitelist(ctx, msg); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.RemoveWhitelist, str,
		),
	}
}

// handleMsgForbidAddr - Handle MsgForbidAddr
func handleMsgForbidAddr(ctx sdk.Context, tk TokenKeeper, msg MsgForbidAddr) (res sdk.Result) {
	if err := tk.ForbidAddress(ctx, msg); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.ForbidAddr {
		str = str + addr.String() + ","
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.Addresses, str,
		),
	}
}

// handleMsgUnForbidAddr - Handle MsgUnForbidAddr
func handleMsgUnForbidAddr(ctx sdk.Context, tk TokenKeeper, msg MsgUnForbidAddr) (res sdk.Result) {
	if err := tk.UnForbidAddress(ctx, msg); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.UnForbidAddr {
		str = str + addr.String() + ","
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Token, msg.Symbol,
			tags.Addresses, str,
		),
	}
}
