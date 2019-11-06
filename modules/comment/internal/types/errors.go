package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceComment sdk.CodespaceType = "comment"

	// 901 ~ 999
	CodeInvalidSymbol      sdk.CodeType = 901
	CodeNegativeDonation   sdk.CodeType = 902
	CodeNoTitle            sdk.CodeType = 903
	CodeInvalidContentType sdk.CodeType = 904
	CodeInvalidContent     sdk.CodeType = 905
	CodeContentTooLarge    sdk.CodeType = 906
	CodeInvalidAttitude    sdk.CodeType = 907
	CodeNegativeReward     sdk.CodeType = 908
	CodeNoSuchAsset        sdk.CodeType = 909
	CodeTitleTooLarge      sdk.CodeType = 910
	CodeNoSuchAccount      sdk.CodeType = 911
	CodeMarshalFailed      sdk.CodeType = 914
)

func ErrInvalidSymbol() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeInvalidSymbol, "Invalid Symbol")
}

func ErrNegativeDonation() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeNegativeDonation, "Donation can not be negative")
}

func ErrNoTitle() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeNoTitle, "No title is provided")
}

func ErrTitleTooLarge() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeTitleTooLarge,
		fmt.Sprintf("Title is larger than %d bytes", MaxTitleSize))
}

func ErrInvalidContentType(t int8) sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeInvalidContentType, fmt.Sprintf("'%d' is not a valid content type", t))
}

func ErrInvalidContent() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeInvalidContent, "Content has invalid format")
}

func ErrContentTooLarge() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeContentTooLarge,
		fmt.Sprintf("Content is larger than %d bytes", MaxContentSize))
}

func ErrInvalidAttitude(a int32) sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeInvalidAttitude, fmt.Sprintf("'%d' is not a valid attitude", a))
}

func ErrNegativeReward() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeNegativeReward, "Reward can not be negative")
}

func ErrNoSuchAsset() sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeNoSuchAsset, "No such asset")
}

func ErrNoSuchAccount(acc string) sdk.Error {
	return sdk.NewError(CodeSpaceComment, CodeNoSuchAccount, fmt.Sprintf("No such account: %s", acc))
}
