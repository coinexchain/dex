package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceBancorlite sdk.CodespaceType = "bancorlite"

	// 1001 ~ 1099
	CodeInvalidSymbol         sdk.CodeType = 1001
	CodeNonPositiveSupply     sdk.CodeType = 1002
	CodeNonPositivePrice      sdk.CodeType = 1003
	CodeNonPositiveAmount     sdk.CodeType = 1004
	CodeTradeAmountIsTooLarge sdk.CodeType = 1005
	CodeBancorAlreadyExists   sdk.CodeType = 1006
	CodeNoSuchToken           sdk.CodeType = 1007
	CodeNonOwnerIsProhibited  sdk.CodeType = 1008
	CodeNoBancorExists        sdk.CodeType = 1009
	CodeOwnerIsProhibited     sdk.CodeType = 1010
	CodeStockInPoolOutofBound sdk.CodeType = 1011
	CodeMoneyCrossLimit       sdk.CodeType = 1012
	CodeUnMarshalFailed       sdk.CodeType = 1013
	CodeMarshalFailed         sdk.CodeType = 1014
)

func ErrInvalidSymbol() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeInvalidSymbol, "Invalid Symbol")
}

func ErrNonPositiveSupply() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonPositiveSupply, "Non-positive supply is invalid")
}

func ErrNonPositivePrice() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonPositivePrice, "Non-positive price is invalid")
}

func ErrNonPositiveAmount() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonPositiveAmount, "Negative or zero amount is invalid")
}

func ErrTradeAmountIsTooLarge() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeTradeAmountIsTooLarge, "Trade amount is too large")
}

func ErrBancorAlreadyExists() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeBancorAlreadyExists, "The Bancor pool is already created")
}

func ErrNoSuchToken() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNoSuchToken, "No such token.")
}

func ErrNonOwnerIsProhibited() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonOwnerIsProhibited, "Non-owner of this token can not create Bancor pool for it.")
}

func ErrNoBancorExists() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNoBancorExists, "The Bancor pool for this token does not exist")
}

func ErrOwnerIsProhibited() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeOwnerIsProhibited, "The token's owner can not trade with the token's Bancor pool.")
}

func ErrStockInPoolOutofBound() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeStockInPoolOutofBound, "The stock in Bancor pool will be out of bound.")
}

func ErrMoneyCrossLimit(moneyErr string) sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeMoneyCrossLimit, "The money amount in this trade is "+moneyErr+" the limited value.")
}
