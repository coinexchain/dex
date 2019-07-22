package asset

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset/internal/types"
)

// NewHandler returns a handler for "asset" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgIssueToken:
			return handleMsgIssueToken(ctx, keeper, msg)
		case types.MsgTransferOwnership:
			return handleMsgTransferOwnership(ctx, keeper, msg)
		case types.MsgMintToken:
			return handleMsgMintToken(ctx, keeper, msg)
		case types.MsgBurnToken:
			return handleMsgBurnToken(ctx, keeper, msg)
		case types.MsgForbidToken:
			return handleMsgForbidToken(ctx, keeper, msg)
		case types.MsgUnForbidToken:
			return handleMsgUnForbidToken(ctx, keeper, msg)
		case types.MsgAddTokenWhitelist:
			return handleMsgAddTokenWhitelist(ctx, keeper, msg)
		case types.MsgRemoveTokenWhitelist:
			return handleMsgRemoveTokenWhitelist(ctx, keeper, msg)
		case types.MsgForbidAddr:
			return handleMsgForbidAddr(ctx, keeper, msg)
		case types.MsgUnForbidAddr:
			return handleMsgUnForbidAddr(ctx, keeper, msg)
		case types.MsgModifyTokenInfo:
			return handleMsgModifyTokenInfo(ctx, keeper, msg)

		default:
			errMsg := "Unrecognized asset Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgIssueToken - Handle MsgIssueToken
func handleMsgIssueToken(ctx sdk.Context, keeper Keeper, msg types.MsgIssueToken) sdk.Result {
	issueFee := keeper.GetParams(ctx).IssueTokenFee
	if len(msg.Symbol) == types.RareSymbolLength {
		issueFee = keeper.GetParams(ctx).IssueRareTokenFee
	}

	if err := keeper.DeductFee(ctx, msg.Owner, issueFee); err != nil {
		return err.Result()
	}

	err := keeper.IssueToken(ctx, msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable, msg.URL, msg.Description)

	if err != nil {
		return err.Result()
	}

	if err := keeper.SendCoinsFromAssetModuleToAccount(ctx, msg.Owner, types.NewTokenCoins(msg.Symbol, msg.TotalSupply)); err != nil {
		return err.Result()
	}

	//if err := keeper.AddToken(ctx, msg.Owner, types.NewTokenCoins(msg.Symbol, msg.TotalSupply)); err != nil {
	//	return err.Result()
	//}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyTokenOwner, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgTransferOwnership - Handle MsgTransferOwnership
func handleMsgTransferOwnership(ctx sdk.Context, keeper Keeper, msg types.MsgTransferOwnership) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.TransferOwnership(ctx, msg.Symbol, msg.OriginalOwner, msg.NewOwner); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyTokenOwner, msg.NewOwner.String()),
			sdk.NewAttribute(types.AttributeKeyOriginalOwner, msg.OriginalOwner.String()),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgMintToken - Handle MsgMintToken
func handleMsgMintToken(ctx sdk.Context, keeper Keeper, msg types.MsgMintToken) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.MintToken(ctx, msg.Symbol, msg.OwnerAddress, msg.Amount); err != nil {
		return err.Result()
	}
	if err := keeper.SendCoinsFromAssetModuleToAccount(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	//if err := keeper.AddToken(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
	//	return err.Result()
	//}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyMintAmount, strconv.FormatInt(msg.Amount, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgBurnToken - Handle MsgBurnToken
func handleMsgBurnToken(ctx sdk.Context, keeper Keeper, msg types.MsgBurnToken) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	if err := keeper.SendCoinsFromAccountToAssetModule(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	if err := keeper.BurnToken(ctx, msg.Symbol, msg.OwnerAddress, msg.Amount); err != nil {
		return err.Result()
	}

	//if err := keeper.SubtractToken(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
	//	return err.Result()
	//}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyMintAmount, strconv.FormatInt(msg.Amount, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgForbidToken - Handle ForbidToken msg
func handleMsgForbidToken(ctx sdk.Context, keeper Keeper, msg types.MsgForbidToken) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.ForbidToken(ctx, msg.Symbol, msg.OwnerAddress); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgUnForbidToken - Handle UnForbidToken msg
func handleMsgUnForbidToken(ctx sdk.Context, keeper Keeper, msg types.MsgUnForbidToken) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.UnForbidToken(ctx, msg.Symbol, msg.OwnerAddress); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgAddTokenWhitelist - Handle AddTokenWhitelist msg
func handleMsgAddTokenWhitelist(ctx sdk.Context, keeper Keeper, msg types.MsgAddTokenWhitelist) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.AddTokenWhitelist(ctx, msg.Symbol, msg.OwnerAddress, msg.Whitelist); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyAddWhitelist, str),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgRemoveTokenWhitelist - Handle RemoveTokenWhitelist msg
func handleMsgRemoveTokenWhitelist(ctx sdk.Context, keeper Keeper, msg types.MsgRemoveTokenWhitelist) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.RemoveTokenWhitelist(ctx, msg.Symbol, msg.OwnerAddress, msg.Whitelist); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyRemoveWhitelist, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgForbidAddr - Handle MsgForbidAddr
func handleMsgForbidAddr(ctx sdk.Context, keeper Keeper, msg types.MsgForbidAddr) (res sdk.Result) {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.ForbidAddress(ctx, msg.Symbol, msg.OwnerAddr, msg.Addresses); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Addresses {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyAddr, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgUnForbidAddr - Handle MsgUnForbidAddr
func handleMsgUnForbidAddr(ctx sdk.Context, keeper Keeper, msg types.MsgUnForbidAddr) (res sdk.Result) {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.UnForbidAddress(ctx, msg.Symbol, msg.OwnerAddr, msg.Addresses); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Addresses {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyAddr, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgModifyTokenInfo - Handle MsgModifyTokenInfo
func handleMsgModifyTokenInfo(ctx sdk.Context, keeper Keeper, msg types.MsgModifyTokenInfo) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.ModifyTokenInfo(ctx, msg.Symbol, msg.OwnerAddress, msg.URL, msg.Description); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeAsset,
			sdk.NewAttribute(types.AttributeKeyToken, msg.Symbol)),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyURL, msg.URL),
		),
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "asset"),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
