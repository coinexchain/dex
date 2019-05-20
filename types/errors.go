package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const CodespaceDEX sdk.CodespaceType = "dex"

// CET error codes
const (
	CodeUnactivatedAddress sdk.CodeType = 1
	CodeMemoMissing        sdk.CodeType = 2
)

func ErrUnactivatedAddress(msg string) sdk.Error {
	return sdk.NewError(CodespaceDEX, CodeUnactivatedAddress, msg)
}

func ErrMemoMissing(msg string) sdk.Error {
	return sdk.NewError(CodespaceDEX, CodeMemoMissing, msg)
}
