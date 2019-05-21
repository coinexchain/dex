package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx = "bankx"

	CodeInsufficientCETForActivatingFee = 19
)

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivatingFee, "Insufficient CET for Activating fees")
}
