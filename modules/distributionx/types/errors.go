package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodespaceDistrx sdk.CodespaceType = "distrx"

	CodeInvalidFromAddr          sdk.CodeType = 801
	CodeInvalidDonation          sdk.CodeType = 802
	CodeMemoRequiredWithdrawAddr sdk.CodeType = 803
)

func ErrorInvalidFromAddr() sdk.Error {
	return sdk.NewError(CodespaceDistrx, CodeInvalidFromAddr, "invalid from address")
}

func ErrorInvalidDonation(format string) sdk.Error {
	return sdk.NewError(CodespaceDistrx, CodeInvalidDonation, format)
}

func ErrMemoRequiredWithdrawAddr(address string) sdk.Error {
	return sdk.NewError(CodespaceDistrx, CodeMemoRequiredWithdrawAddr, fmt.Sprintf("cannot set memo-required address %s be withdraw address", address))
}
