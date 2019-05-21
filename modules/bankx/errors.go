package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceBankx = "bankx"

	CodeInsufficientCETForActivatedFee = 19
	CodeInvalidActivatedFee            = 20
)

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivatedFee, "Insufficient CET for Activating fees")
}
