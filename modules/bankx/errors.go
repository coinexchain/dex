package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx sdk.CodespaceType = "bankx"

	CodeUnactivatedAddress             sdk.CodeType = 111
	CodeMemoMissing                    sdk.CodeType = 112
	CodeInsufficientCETForActivatedFee sdk.CodeType = 113
	CodeInvalidActivatedFee            sdk.CodeType = 114
)

func ErrUnactivatedAddress(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeUnactivatedAddress, msg)
}

func ErrMemoMissing() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeMemoMissing, "memo is empty")
}

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivatedFee, "Insufficient CET for Activating fees")
}
