package asset

import (
	asset_types "github.com/coinexchain/dex/modules/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAsset sdk.CodespaceType = asset_types.ModuleName

	// 501 ~ 599
	CodeInvalidTokenName        sdk.CodeType = 501
	CodeInvalidTokenSymbol      sdk.CodeType = 502
	CodeInvalidTokenSupply      sdk.CodeType = 503
	CodeInvalidTokenOwner       sdk.CodeType = 504
	CodeTokenNotFound           sdk.CodeType = 505
	CodeInvalidTotalMint        sdk.CodeType = 506
	CodeInvalidTotalBurn        sdk.CodeType = 507
	CodeDuplicateTokenSymbol    sdk.CodeType = 508
	CodeInvalidTokenForbidden   sdk.CodeType = 509
	CodeInvalidTokenWhitelist   sdk.CodeType = 510
	CodeInvalidAddress          sdk.CodeType = 511
	CodeInvalidTokenURL         sdk.CodeType = 512
	CodeInvalidTokenDescription sdk.CodeType = 513
)

func ErrorInvalidTokenName(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenName, fmt)
}
func ErrorInvalidTokenSymbol(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenSymbol, fmt)
}
func ErrorInvalidTokenSupply(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenSupply, fmt)
}
func ErrorInvalidTokenOwner(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenOwner, fmt)
}
func ErrorTokenNotFound(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeTokenNotFound, fmt)
}
func ErrorInvalidTokenMint(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTotalMint, fmt)
}
func ErrorInvalidTokenBurn(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTotalBurn, fmt)
}
func ErrorDuplicateTokenSymbol(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeDuplicateTokenSymbol, fmt)
}
func ErrorInvalidTokenForbidden(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenForbidden, fmt)
}
func ErrorInvalidTokenWhitelist(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenWhitelist, fmt)
}
func ErrorInvalidAddress(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidAddress, fmt)
}
func ErrorInvalidTokenURL(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenURL, fmt)
}
func ErrorInvalidTokenDescription(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenDescription, fmt)
}
