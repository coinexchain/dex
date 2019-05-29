package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx sdk.CodespaceType = "bankx"

	CodeUnactivatedAddress              sdk.CodeType = 111
	CodeMemoMissing                     sdk.CodeType = 112
	CodeInsufficientCETForActivationFee sdk.CodeType = 113
	CodeInvalidActivationFee            sdk.CodeType = 114
	CodeInvalidUnlockTime               sdk.CodeType = 115
	CodeCetCantBeLocked                 sdk.CodeType = 116
)

func ErrUnactivatedAddress(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeUnactivatedAddress, msg)
}

func ErrMemoMissing() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeMemoMissing, "memo is empty")
}

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivationFee, "Insufficient CET for Activating fees")
}

func ErrUnlockTime(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidUnlockTime, msg)
}

func ErrCetCantBeLocked(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeCetCantBeLocked, msg)
}
