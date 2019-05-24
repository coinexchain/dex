package asset

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceAsset = ModuleName

	CodeInvalidTokenName     = 201
	CodeInvalidTokenSymbol   = 202
	CodeInvalidTokenSupply   = 203
	CodeInvalidTokenOwner    = 204
	CodeNoTokenPersist       = 205
	CodeInvalidTotalMint     = 206
	CodeInvalidTotalBurn     = 207
	CodeDuplicateTokenSymbol = 208
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
func ErrorNoTokenPersist(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeNoTokenPersist, fmt)
}
func ErrorInvalidTotalMint(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTotalMint, fmt)
}
func ErrorInvalidTotalBurn(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTotalBurn, fmt)
}
func ErrorDuplicateTokenSymbol(fmt string) sdk.Error {
	return sdk.NewError(CodeSpaceAsset, CodeDuplicateTokenSymbol, fmt)
}
