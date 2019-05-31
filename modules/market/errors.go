package market

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceMarket sdk.CodespaceType = "market"

	CodeInvalidToken          = 120
	CodeInvalidPricePrecision = 121
	CodeInvalidTokenIssuer    = 122
	CodeSendTokenFailed       = 123
	CodeNoStoreEngine         = 124
	CodeInvalidAddress        = 125
	CodeNotExistKeyInStore    = 126
	CodeNotHaveSufficientCoin = 127
	CodeInvalidTradeSide      = 128
	CodeInvalidOrderType      = 129
	CodeInvalidSymbol         = 130
	CodeTokenFrozenByIssuer   = 131
	CodeInvalidOrderID        = 132
)

func ErrTokenNoExist() sdk.Error {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidToken, "Token not exsit")
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

func ErrNoHaveSufficientCoins() sdk.Error {

	return sdk.NewError(CodeSpaceMarket, CodeNotHaveSufficientCoin, "Not sufficient coin")
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

func ErrTokenFrozenByIssuer() sdk.Error {

	return sdk.NewError(CodeSpaceMarket, CodeTokenFrozenByIssuer, "Token is frozen by the issuer")
}
