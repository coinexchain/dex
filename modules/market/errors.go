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
)

func ErrTokenNoExsit() sdk.Result {

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
