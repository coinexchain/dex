package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceBancorlite sdk.CodespaceType = "bancorlite"

	// 1001 ~ 1099
	CodeInvalidSymbol                sdk.CodeType = 1001
	CodeNonPositiveSupply            sdk.CodeType = 1002
	CodeNonPositivePrice             sdk.CodeType = 1003
	CodeNonPositiveAmount            sdk.CodeType = 1004
	CodeTradeAmountIsTooLarge        sdk.CodeType = 1005
	CodeBancorAlreadyExists          sdk.CodeType = 1006
	CodeNoSuchToken                  sdk.CodeType = 1007
	CodeNonOwnerIsProhibited         sdk.CodeType = 1008
	CodeNoBancorExists               sdk.CodeType = 1009
	CodeOwnerIsProhibited            sdk.CodeType = 1010
	CodeStockInPoolOutOfBound        sdk.CodeType = 1011
	CodeMoneyCrossLimit              sdk.CodeType = 1012
	CodeUnMarshalFailed              sdk.CodeType = 1013
	CodeMarshalFailed                sdk.CodeType = 1014
	CodeNegativeInitPrice            sdk.CodeType = 1015
	CodeNonMarketExist               sdk.CodeType = 1016
	CodeNotBancorOwner               sdk.CodeType = 1017
	CodeCancelTimeNotArrived         sdk.CodeType = 1018
	CodeGetMarketExePriceFailed      sdk.CodeType = 1019
	CodeInitPriceBigThanMaxPrice     sdk.CodeType = 1020
	CodeCancelEnableTimeNegative     sdk.CodeType = 1021
	CodeTradeQuantityTooSmall        sdk.CodeType = 1022
	CodeTokenForbiddenByOwner        sdk.CodeType = 1023
	CodeStockSupplyPrecisionNotMatch sdk.CodeType = 1024
	CodeErrPriceFmt                  sdk.CodeType = 1025
	CodeStockAmountPrecisionNotMatch sdk.CodeType = 1026
)

func ErrInvalidSymbol() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeInvalidSymbol, "Invalid Symbol")
}

func ErrNonPositiveSupply() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonPositiveSupply, "Non-positive supply is invalid")
}

func ErrPriceFmt() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeErrPriceFmt, "Invalid Price format")
}

func ErrNonPositivePrice() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonPositivePrice, "Non-positive price is invalid")
}

func ErrNegativePrice() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNegativeInitPrice, "Negative init price is invalid")
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

func ErrNonMarketExist() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNonMarketExist, "No corresponding market exist")
}

func ErrNoBancorExists() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNoBancorExists, "The Bancor pool for this token does not exist")
}

func ErrOwnerIsProhibited() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeOwnerIsProhibited, "The token's owner can not trade with the token's Bancor pool.")
}

func ErrNotBancorOwner() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeNotBancorOwner, "The sender is not the bancor owner")
}

func ErrEarliestCancelTimeNotArrive() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeCancelTimeNotArrived, "The time when bancor can be canceled has not arrived")
}

func ErrStockInPoolOutofBound() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeStockInPoolOutOfBound, "The stock in Bancor pool will be out of bound.")
}

func ErrMoneyCrossLimit(moneyErr string) sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeMoneyCrossLimit, "The money amount in this trade is "+moneyErr+" the limited value.")
}

func ErrGetMarketPrice(err string) sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeGetMarketExePriceFailed, err)
}

func ErrTradeQuantityTooSmall(amount int64) sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeTradeQuantityTooSmall, "The trade commission (%d) too small", amount)
}

func ErrPriceConfiguration() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeInitPriceBigThanMaxPrice, "The init price is bigger than max price")
}

func ErrEarliestCancelTimeIsNegative() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeCancelEnableTimeNegative, "The cancellation enable time is negative")
}

func ErrTokenForbiddenByOwner() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeTokenForbiddenByOwner, "token is forbidden by its owner")
}

func ErrStockSupplyPrecisionNotMatch() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeStockSupplyPrecisionNotMatch, "stock supply not match the stock precision")
}

func ErrStockAmountPrecisionNotMatch() sdk.Error {
	return sdk.NewError(CodeSpaceBancorlite, CodeStockAmountPrecisionNotMatch, "stock amount not match the stock precision")
}
