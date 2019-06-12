package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx sdk.CodespaceType = "bankx"

	CodeMemoMissing                     sdk.CodeType = 111
	CodeInsufficientCETForActivationFee sdk.CodeType = 112
	CodeInvalidActivationFee            sdk.CodeType = 113
	CodeInvalidUnlockTime               sdk.CodeType = 114
	CodeTokenForbiddenByOwner           sdk.CodeType = 115
	CodeInvalidLockCoinsFee             sdk.CodeType = 116
)

func ErrMemoMissing() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeMemoMissing, "memo is empty")
}

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivationFee, "Insufficient CET for Activating fees")
}

func ErrUnlockTime(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidUnlockTime, msg)
}

func ErrTokenForbiddenByOwner(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeTokenForbiddenByOwner, msg)
}
