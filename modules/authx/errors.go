package authx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceAuthX sdk.CodespaceType = "authx"

	CodeAccountxInvalid    sdk.CodeType = 200
	CodeInvalidMinGasPrice sdk.CodeType = 201
)

func ErrInvalidAccountx(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeAccountxInvalid, msg)
}

func ErrInvalidMinGasPrice(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeInvalidMinGasPrice, msg)
}
