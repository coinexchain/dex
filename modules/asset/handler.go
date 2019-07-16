package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset/types"
)

// NewHandler returns a handler for "asset" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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
		case types.MsgModifyTokenURL:
			return handleMsgModifyTokenURL(ctx, keeper, msg)
		case types.MsgModifyTokenDescription:
			return handleMsgModifyTokenDescription(ctx, keeper, msg)

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

	if err := keeper.AddToken(ctx, msg.Owner, newTokenCoins(msg.Symbol, msg.TotalSupply)); err != nil {
		return err.Result()
	}

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Owner, msg.Owner.String(),
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.OriginalOwner, msg.OriginalOwner.String(),
		//	tags.NewOwner, msg.NewOwner.String(),
		//),
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

	if err := keeper.AddToken(ctx, msg.OwnerAddress, newTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Amt, strconv.FormatInt(msg.Amount, 10),
		//),
	}
}

// handleMsgBurnToken - Handle MsgBurnToken
func handleMsgBurnToken(ctx sdk.Context, keeper Keeper, msg types.MsgBurnToken) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.BurnToken(ctx, msg.Symbol, msg.OwnerAddress, msg.Amount); err != nil {
		return err.Result()
	}

	if err := keeper.SubtractToken(ctx, msg.OwnerAddress, newTokenCoins(msg.Symbol, msg.Amount)); err != nil {
		return err.Result()
	}

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Amt, strconv.FormatInt(msg.Amount, 10),
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.AddWhitelist, str,
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.RemoveWhitelist, str,
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Addresses, str,
		//),
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

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Addresses, str,
		//),
	}
}

// handleMsgModifyTokenURL - Handle MsgModifyTokenURL
func handleMsgModifyTokenURL(ctx sdk.Context, keeper Keeper, msg types.MsgModifyTokenURL) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.ModifyTokenURL(ctx, msg.Symbol, msg.OwnerAddress, msg.URL); err != nil {
		return err.Result()
	}

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.URL, msg.URL,
		//),
	}
}

// handleMsgModifyTokenDescription - Handle MsgModifyTokenDescription
func handleMsgModifyTokenDescription(ctx sdk.Context, keeper Keeper, msg types.MsgModifyTokenDescription) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	if err := keeper.ModifyTokenDescription(ctx, msg.Symbol, msg.OwnerAddress, msg.Description); err != nil {
		return err.Result()
	}

	return sdk.Result{
		//Tags: sdk.NewTags(
		//	tags.Category, tags.TxCategory,
		//	tags.Token, msg.Symbol,
		//	tags.Description, msg.Description,
		//),
	}
}

func newTokenCoins(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(denom, amount))
}
