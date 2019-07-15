package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx sdk.CodespaceType = ModuleName

	// 301 ～ 399
	CodeMemoMissing                     sdk.CodeType = 301
	CodeInsufficientCETForActivationFee sdk.CodeType = 302
	CodeInvalidActivationFee            sdk.CodeType = 303
	CodeInvalidUnlockTime               sdk.CodeType = 304
	CodeTokenForbiddenByOwner           sdk.CodeType = 305
	CodeInvalidLockCoinsFee             sdk.CodeType = 306
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
