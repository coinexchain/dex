package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceStakingX sdk.CodespaceType = "stakingx"

	CodeMinSelfDelegationBelowRequired sdk.CodeType = 201
)

func ErrMinSelfDelegationBelowRequired(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceStakingX, CodeMinSelfDelegationBelowRequired, msg)
}
