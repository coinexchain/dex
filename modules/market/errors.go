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
)

func ErrTokenNoExist() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidToken, "Token not exsit").Result()
}

func ErrInvalidPricePrecision() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidPricePrecision, "Price precision out of range").Result()
}

func ErrInvalidTokenIssuer() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidTokenIssuer, "Invalid token issuer").Result()
}

func ErrSendTokenFailed(errStr string) sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeSendTokenFailed, "Send token failed %s", errStr).Result()
}

func ErrNoStoreEngine() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeNoStoreEngine, "market No store engine").Result()
}

func ErrInvalidAddress() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidAddress, "Invalid address").Result()
}

func ErrNoExistKeyInStore() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeNotExistKeyInStore, "Not exist key in store").Result()
}

func ErrNoHaveSufficientCoins() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeNotHaveSufficientCoin, "Not sufficient coin").Result()
}

func ErrInvalidTradeSide() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidTradeSide, "Invalid trade side").Result()
}

func ErrInvalidOrderType() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderType, "Invalid order type").Result()
}

func ErrInvalidSymbol() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "Invalid trade symbol").Result()
}

func ErrTokenFrozenByIssuer() sdk.Result {

	return sdk.NewError(CodeSpaceMarket, CodeTokenFrozenByIssuer, "Token is frozen by the issuer").Result()
}
