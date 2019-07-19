package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceMarket sdk.CodespaceType = "market"

	// 601 ~ 699
	CodeInvalidToken          sdk.CodeType = 601
	CodeInvalidPricePrecision sdk.CodeType = 602
	CodeInvalidTokenIssuer    sdk.CodeType = 603
	CodeInvalidAddress        sdk.CodeType = 606
	CodeNotExistKeyInStore    sdk.CodeType = 607
	CodeInsufficientCoin      sdk.CodeType = 608
	CodeInvalidTradeSide      sdk.CodeType = 609
	CodeInvalidOrderType      sdk.CodeType = 610
	CodeInvalidSymbol         sdk.CodeType = 611
	CodeTokenForbidByIssuer   sdk.CodeType = 612
	CodeInvalidOrderID        sdk.CodeType = 613
	CodeMarshalFailed         sdk.CodeType = 614
	CodeUnMarshalFailed       sdk.CodeType = 615
	CodeNotFindOrder          sdk.CodeType = 616
	CodeNotMatchSender        sdk.CodeType = 617
	CodeInvalidTime           sdk.CodeType = 618
	CodeAddressForbidByIssuer sdk.CodeType = 619
	CodeOrderQuantityToSmall  sdk.CodeType = 620
	CodeStockNoHaveCetTrade   sdk.CodeType = 621
	CodeRepeatTrade           sdk.CodeType = 622
	CodeNotAllowedOffline     sdk.CodeType = 623
	CodeInvalidOrderAmount    sdk.CodeType = 624
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

func ErrInvalidPrice(price int64) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidPricePrecision, "Price out of range [0, 9E18], actual price :  ", price)
}

func ErrInvalidTokenIssuer() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidTokenIssuer, "Invalid token issuer")
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
