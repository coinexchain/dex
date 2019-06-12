package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAuthX sdk.CodespaceType = "authx"

	CodeInvalidMinGasPriceLimit sdk.CodeType = 201
)

func ErrInvalidMinGasPriceLimit(limit int64) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeInvalidMinGasPriceLimit,
		"invalid minimum gas price limit: %d", limit)
}
