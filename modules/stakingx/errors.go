package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceStakingX sdk.CodespaceType = "stakingx"

	// 401 ~ 499
	CodeInvalidMinSelfDelegation        sdk.CodeType = 401
	CodeMinSelfDelegationBelowRequired  sdk.CodeType = 402
	CodeBelowMinMandatoryCommissionRate sdk.CodeType = 403
)

func ErrInvalidMinSelfDelegation(val sdk.Int) sdk.Error {
	return sdk.NewError(CodeSpaceStakingX, CodeInvalidMinSelfDelegation,
		"invalid min gas price: %v", val)
}

func ErrMinSelfDelegationBelowRequired(expected, actual sdk.Int) sdk.Error {
	return sdk.NewError(CodeSpaceStakingX, CodeMinSelfDelegationBelowRequired,
		"minimum self-delegation is %v, less than %v", actual, expected)
}

func ErrRateBelowMinMandatoryCommissionRate(expected, actual sdk.Dec) sdk.Error {
	return sdk.NewError(CodeSpaceStakingX, CodeBelowMinMandatoryCommissionRate,
		"commission rate is %v, less than min mandatory commission rate %v", actual, expected)
}
