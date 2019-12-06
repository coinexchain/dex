package asset

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
)

// NewHandler returns a handler for "asset" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

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
			return dex.ErrUnknownRequest(ModuleName, msg)
		}
	}
}

// handleMsgIssueToken - Handle MsgIssueToken
func handleMsgIssueToken(ctx sdk.Context, keeper Keeper, msg types.MsgIssueToken) sdk.Result {
	issueFee := keeper.GetParams(ctx).GetIssueTokenFee(msg.Symbol)

	if err := keeper.DeductIssueFee(ctx, msg.Owner, issueFee); err != nil {
		return err.Result()
	}

	err := keeper.IssueToken(ctx, msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable, msg.URL, msg.Description, msg.Identity)

	if err != nil {
		return err.Result()
	}

	if err := keeper.SendCoinsFromAssetModuleToAccount(ctx, msg.Owner, types.NewTokenCoins(msg.Symbol, msg.TotalSupply)); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyTokenOwner, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgTransferOwnership - Handle MsgTransferOwnership
func handleMsgTransferOwnership(ctx sdk.Context, keeper Keeper, msg types.MsgTransferOwnership) sdk.Result {
	if err := keeper.TransferOwnership(ctx, msg.Symbol, msg.OriginalOwner, msg.NewOwner); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OriginalOwner.String()),
		),
		sdk.NewEvent(
			types.EventTypeTransferOwnership,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
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
	if err := keeper.MintToken(ctx, msg.Symbol, msg.OwnerAddress, msg.Amount); err != nil {
		return err.Result()
	}
	if err := keeper.SendCoinsFromAssetModuleToAccount(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgBurnToken - Handle MsgBurnToken
func handleMsgBurnToken(ctx sdk.Context, keeper Keeper, msg types.MsgBurnToken) sdk.Result {
	if err := keeper.SendCoinsFromAccountToAssetModule(ctx, msg.OwnerAddress, types.NewTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	if err := keeper.BurnToken(ctx, msg.Symbol, msg.OwnerAddress, msg.Amount); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgForbidToken - Handle ForbidToken msg
func handleMsgForbidToken(ctx sdk.Context, keeper Keeper, msg types.MsgForbidToken) sdk.Result {
	if err := keeper.ForbidToken(ctx, msg.Symbol, msg.OwnerAddress); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeForbidToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgUnForbidToken - Handle UnForbidToken msg
func handleMsgUnForbidToken(ctx sdk.Context, keeper Keeper, msg types.MsgUnForbidToken) sdk.Result {
	if err := keeper.UnForbidToken(ctx, msg.Symbol, msg.OwnerAddress); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeUnForbidToken,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgAddTokenWhitelist - Handle AddTokenWhitelist msg
func handleMsgAddTokenWhitelist(ctx sdk.Context, keeper Keeper, msg types.MsgAddTokenWhitelist) sdk.Result {
	if err := keeper.AddTokenWhitelist(ctx, msg.Symbol, msg.OwnerAddress, msg.Whitelist); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeAddTokenWhitelist,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddrList, str),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgRemoveTokenWhitelist - Handle RemoveTokenWhitelist msg
func handleMsgRemoveTokenWhitelist(ctx sdk.Context, keeper Keeper, msg types.MsgRemoveTokenWhitelist) sdk.Result {
	if err := keeper.RemoveTokenWhitelist(ctx, msg.Symbol, msg.OwnerAddress, msg.Whitelist); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Whitelist {
		str = str + addr.String()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeRemoveTokenWhitelist,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddrList, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgForbidAddr - Handle MsgForbidAddr
func handleMsgForbidAddr(ctx sdk.Context, keeper Keeper, msg types.MsgForbidAddr) (res sdk.Result) {
	if err := keeper.ForbidAddress(ctx, msg.Symbol, msg.OwnerAddr, msg.Addresses); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Addresses {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddr.String()),
		),
		sdk.NewEvent(types.EventTypeForbidAddr,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddrList, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgUnForbidAddr - Handle MsgUnForbidAddr
func handleMsgUnForbidAddr(ctx sdk.Context, keeper Keeper, msg types.MsgUnForbidAddr) (res sdk.Result) {
	if err := keeper.UnForbidAddress(ctx, msg.Symbol, msg.OwnerAddr, msg.Addresses); err != nil {
		return err.Result()
	}

	var str string
	for _, addr := range msg.Addresses {
		str = str + addr.String() + ","
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddr.String()),
		),
		sdk.NewEvent(types.EventTypeUnForbidAddr,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddrList, str),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgModifyTokenInfo - Handle MsgModifyTokenInfo
func handleMsgModifyTokenInfo(ctx sdk.Context, keeper Keeper, msg types.MsgModifyTokenInfo) sdk.Result {
	token := keeper.GetToken(ctx, msg.Symbol)
	if token == nil {
		return types.ErrTokenNotFound(msg.Symbol).Result()
	}

	newURL, newDesc, newID, newName, newSupply,
		newMintable, newBurnable, newAddrForbiddable, newTokenForbiddable,
		err := CollectTokenModificationInfo(token, msg)
	if err != nil {
		return err.Result()
	}

	if err := keeper.ModifyTokenInfo(ctx, msg.Symbol, msg.OwnerAddress,
		newURL, newDesc, newID, newName, newSupply,
		newMintable, newBurnable, newAddrForbiddable, newTokenForbiddable); err != nil {

		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress.String()),
		),
		sdk.NewEvent(types.EventTypeModifyTokenInfo,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyURL, msg.URL),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
			sdk.NewAttribute(types.AttributeKeyIdentity, msg.Identity),
		),
	})
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func CollectTokenModificationInfo(token types.Token, msg types.MsgModifyTokenInfo) (
	newURL, newDesc, newID, newName string, newSupply sdk.Int,
	newMintable, newBurnable, newAddrForbiddable, newTokenForbiddable bool,
	sdkErr sdk.Error) {

	var err error

	// modifiable fields
	newURL = getNewStringVal(msg.URL, token.GetURL())
	newDesc = getNewStringVal(msg.Description, token.GetDescription())
	newID = getNewStringVal(msg.Identity, token.GetIdentity())
	newName = getNewStringVal(msg.Name, token.GetName())
	if newSupply, err = getNewIntVal(msg.TotalSupply, token.GetTotalSupply()); err != nil {
		sdkErr = types.ErrInvalidTokenInfo("TotalSupply", msg.TotalSupply)
		return
	}
	if newMintable, err = getNewBoolVal(msg.Mintable, token.GetMintable()); err != nil {
		sdkErr = types.ErrInvalidTokenInfo("Mintable", msg.Mintable)
		return
	}
	if newBurnable, err = getNewBoolVal(msg.Burnable, token.GetBurnable()); err != nil {
		sdkErr = types.ErrInvalidTokenInfo("Burnable", msg.Burnable)
		return
	}
	if newAddrForbiddable, err = getNewBoolVal(msg.AddrForbiddable, token.GetAddrForbiddable()); err != nil {
		sdkErr = types.ErrInvalidTokenInfo("AddrForbiddable", msg.AddrForbiddable)
		return
	}
	if newTokenForbiddable, err = getNewBoolVal(msg.TokenForbiddable, token.GetTokenForbiddable()); err != nil {
		sdkErr = types.ErrInvalidTokenInfo("TokenForbiddable", msg.TokenForbiddable)
		return
	}
	return
}

func getNewStringVal(newVal, oldVal string) string {
	if newVal == types.DoNotModifyTokenInfo {
		return oldVal
	}
	return newVal
}
func getNewBoolVal(newVal string, oldVal bool) (bool, error) {
	if newVal == types.DoNotModifyTokenInfo {
		return oldVal, nil
	}
	return strconv.ParseBool(newVal)
}
func getNewIntVal(newVal string, oldVal sdk.Int) (sdk.Int, error) {
	if newVal == types.DoNotModifyTokenInfo {
		return oldVal, nil
	}
	n, ok := sdk.NewIntFromString(newVal)
	if !ok {
		return oldVal, errors.New(newVal)
	}
	return n, nil
}
