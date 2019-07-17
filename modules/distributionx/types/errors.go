package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodespaceDistrx sdk.CodespaceType = "distrx"

	CodeInvalidFromAddr sdk.CodeType = 801
	CodeInvalidDonation sdk.CodeType = 802
)

func ErrorInvalidFromAddr() sdk.Error {
	return sdk.NewError(CodespaceDistrx, CodeInvalidFromAddr, "invalid from address")
}

func ErrorInvalidDonation(format string) sdk.Error {
	return sdk.NewError(CodespaceDistrx, CodeInvalidDonation, format)
}
