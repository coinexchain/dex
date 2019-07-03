package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAuthX sdk.CodespaceType = "authx"

	// 201 ～ 299
	CodeInvalidMinGasPriceLimit sdk.CodeType = 201
	CodeGasPriceTooLow          sdk.CodeType = 202
)

func ErrInvalidMinGasPriceLimit(limit sdk.Dec) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeInvalidMinGasPriceLimit,
		"invalid minimum gas price limit: %s", limit)
}

func ErrGasPriceTooLow(required, actual sdk.Dec) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeGasPriceTooLow,
		"gas price too low: %s < %s", actual, required)
}
