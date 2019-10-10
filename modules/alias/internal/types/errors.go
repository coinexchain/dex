package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAlias sdk.CodespaceType = "alias"

	// 1101 ~ 1199
	CodeEmptyAlias              sdk.CodeType = 1101
	CodeInvalidAlias            sdk.CodeType = 1102
	CodeUnknownOperation        sdk.CodeType = 1103
	CodeMarshalFailed           sdk.CodeType = 1104
	CodeUnMarshalFailed         sdk.CodeType = 1105
	CodeAliasAlreadyExists      sdk.CodeType = 1106
	CodeNoSuchAlias             sdk.CodeType = 1107
	CodeCanOnlyBeUsedByCetOwner sdk.CodeType = 1108
	CodeMaxAliasCountReached    sdk.CodeType = 1109
)

func ErrEmptyAlias() sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeEmptyAlias, "Empty Alias")
}

func ErrInvalidAlias() sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeInvalidAlias, "Invalid Alias")
}

func ErrAliasAlreadyExists() sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeAliasAlreadyExists, "This alias ready exists in map table")
}

func ErrNoSuchAlias() sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeNoSuchAlias, "No such alias exists")
}

func ErrCanOnlyBeUsedByCetOwner(a string) sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeCanOnlyBeUsedByCetOwner, fmt.Sprintf("This alias '%s' can only be used by CET's owner", a))
}

func ErrMaxAliasCountReached() sdk.Error {
	return sdk.NewError(CodeSpaceAlias, CodeMaxAliasCountReached, "Have reached the maximum alias count and can not add new aliases")
}
