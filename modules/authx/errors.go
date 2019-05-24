package authx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceAuthX sdk.CodespaceType = "authx"

	CodeAccoutxInvalid     sdk.CodeType = 200
	CodeInvalidMinGasPrice sdk.CodeType = 201
)

func ErrInvalidAccoutx(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeAccoutxInvalid, msg)
}

func ErrInvalidMinGasPrice(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeInvalidMinGasPrice, msg)
}
