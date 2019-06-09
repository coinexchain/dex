package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceMarket sdk.CodespaceType = "market"

	CodeInvalidToken          sdk.CodeType = 120
	CodeInvalidPricePrecision sdk.CodeType = 121
	CodeInvalidTokenIssuer    sdk.CodeType = 122
	CodeSendTokenFailed       sdk.CodeType = 123
	CodeNoStoreEngine         sdk.CodeType = 124
	CodeInvalidAddress        sdk.CodeType = 125
	CodeNotExistKeyInStore    sdk.CodeType = 126
	CodeInsufficientCoin      sdk.CodeType = 127
	CodeInvalidTradeSide      sdk.CodeType = 128
	CodeInvalidOrderType      sdk.CodeType = 129
	CodeInvalidSymbol         sdk.CodeType = 130
	CodeTokenForbidByIssuer   sdk.CodeType = 131
	CodeInvalidOrderID        sdk.CodeType = 132
	CodeMarshalFailed         sdk.CodeType = 133
	CodeUnMarshalFailed       sdk.CodeType = 134
	CodeNotFindOrder          sdk.CodeType = 135
	CodeNotMatchSender        sdk.CodeType = 135
	CodeInvalidHeight         sdk.CodeType = 136
	CodeAddressForbidByIssuer sdk.CodeType = 137
	CodeOrderQuantityToSmall  sdk.CodeType = 138
)

func ErrTokenNoExist() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidToken, "Token not exist")
}

func ErrInvalidOrderID() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderID, "Invalid order id")
}

func ErrInvalidPricePrecision() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidPricePrecision, "Price precision out of range")
}

func ErrInvalidPrice() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidPricePrecision, "Price out of range [0, 9E18]")
}

func ErrInvalidTokenIssuer() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidTokenIssuer, "Invalid token issuer")
}

func ErrSendTokenFailed(errStr string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeSendTokenFailed, "Send token failed %s", errStr)
}

func ErrNoStoreEngine() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeNoStoreEngine, "market No store engine")
}

func ErrInvalidAddress() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidAddress, "Invalid address")
}

func ErrNoExistKeyInStore() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeNotExistKeyInStore, "Not exist key in store")
}

func ErrInsufficientCoins() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInsufficientCoin, "Insufficient coin")
}

func ErrInvalidTradeSide() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidTradeSide, "Invalid trade side")
}

func ErrInvalidOrderType() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderType, "Invalid order type")
}

func ErrInvalidSymbol() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "Invalid trade symbol")
}

func ErrTokenForbidByIssuer() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeTokenForbidByIssuer, "Token is frozen by the issuer")
}

func ErrOrderQuantityToSmall() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeOrderQuantityToSmall, "the order's quantity is too small")
}
