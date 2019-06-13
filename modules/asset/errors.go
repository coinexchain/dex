package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAsset sdk.CodespaceType = ModuleName

	CodeInvalidTokenName      sdk.CodeType = 201
	CodeInvalidTokenSymbol    sdk.CodeType = 202
	CodeInvalidTokenSupply    sdk.CodeType = 203
	CodeInvalidTokenOwner     sdk.CodeType = 204
	CodeTokenNotFound         sdk.CodeType = 205
	CodeInvalidTotalMint      sdk.CodeType = 206
	CodeInvalidTotalBurn      sdk.CodeType = 207
	CodeDuplicateTokenSymbol  sdk.CodeType = 208
	CodeInvalidTokenForbidden sdk.CodeType = 209
	CodeInvalidTokenWhitelist sdk.CodeType = 210
	CodeInvalidAddress        sdk.CodeType = 211
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
