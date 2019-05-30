package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx sdk.CodespaceType = "bankx"

	CodeMemoMissing                     sdk.CodeType = 111
	CodeInsufficientCETForActivationFee sdk.CodeType = 112
	CodeInvalidActivationFee            sdk.CodeType = 113
	CodeInvalidUnlockTime               sdk.CodeType = 114
	CodeCetCantBeLocked                 sdk.CodeType = 115
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

func ErrCetCantBeLocked(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeCetCantBeLocked, msg)
}
