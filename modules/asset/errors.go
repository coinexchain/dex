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

func ErrorInvalidTokenName(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenName, fmt)
}
func ErrorInvalidTokenSymbol(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenSymbol, fmt)
}
func ErrorInvalidTokenSupply(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenSupply, fmt)
}
func ErrorInvalidTokenOwner(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenOwner, fmt)
}
func ErrorNoTokenPersist(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeNoTokenPersist, fmt)
}
func ErrorInvalidTotalMint(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTotalMint, fmt)
}
func ErrorInvalidTotalBurn(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTotalBurn, fmt)
}
func ErrorDuplicateTokenSymbol(codespace sdk.CodespaceType, fmt string) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateTokenSymbol, fmt)
}
