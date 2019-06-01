package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceAsset = ModuleName

	CodeInvalidTokenName      = 201
	CodeInvalidTokenSymbol    = 202
	CodeInvalidTokenSupply    = 203
	CodeInvalidTokenOwner     = 204
	CodeTokenNotFound         = 205
	CodeInvalidTotalMint      = 206
	CodeInvalidTotalBurn      = 207
	CodeDuplicateTokenSymbol  = 208
	CodeInvalidForbiddenState = 209
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
func ErrorInvalidForbiddenState(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidForbiddenState, fmt)
}
