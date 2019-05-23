package authx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceAuthX sdk.CodespaceType = "authx"

	CodeAccoutxInvalid sdk.CodeType = 200
)

func ErrInvalidAccoutx(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceAuthX, CodeAccoutxInvalid, msg)
}
