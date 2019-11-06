package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceMarket sdk.CodespaceType = "market"

	// 601 ~ 699
	CodeInvalidToken           sdk.CodeType = 601
	CodeInvalidPricePrecision  sdk.CodeType = 602
	CodeInvalidTokenIssuer     sdk.CodeType = 603
	CodeInvalidPrice           sdk.CodeType = 604
	CodeInvalidAddress         sdk.CodeType = 606
	CodeNotExistKeyInStore     sdk.CodeType = 607
	CodeInsufficientCoin       sdk.CodeType = 608
	CodeInvalidTradeSide       sdk.CodeType = 609
	CodeInvalidOrderType       sdk.CodeType = 610
	CodeInvalidSymbol          sdk.CodeType = 611
	CodeTokenForbidByIssuer    sdk.CodeType = 612
	CodeInvalidOrderID         sdk.CodeType = 613
	CodeMarshalFailed          sdk.CodeType = 614
	CodeOrderNotFound          sdk.CodeType = 616
	CodeNotMatchSender         sdk.CodeType = 617
	CodeInvalidCancelTime      sdk.CodeType = 618
	CodeAddressForbidByIssuer  sdk.CodeType = 619
	CodeInvalidOrderCommission sdk.CodeType = 620
	CodeNotListedAgainstCet    sdk.CodeType = 621
	CodeRepeatTradingPair      sdk.CodeType = 622
	CodeDelistNotAllowed       sdk.CodeType = 623
	CodeInvalidOrderAmount     sdk.CodeType = 624
	CodeInvalidExistBlocks     sdk.CodeType = 626
	CodeInvalidTimeInForce     sdk.CodeType = 627
	CodeOrderAlreadyExist      sdk.CodeType = 630
	CodeDelistRequestExist     sdk.CodeType = 632
	CodeInvalidMarket          sdk.CodeType = 633
)

func ErrFailedParseParam() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "Failed to parse param")
}

func ErrFailedMarshal() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "Marshal failed")
}

func ErrInvalidExistBlocks(eb int64) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidExistBlocks, fmt.Sprintf("Invalid existence time : %d; The range of expected values [0, +∞] ", eb))
}

func ErrInvalidTimeInForce(tif int64) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidTimeInForce, fmt.Sprintf("Invalid timeInForce : %d; The valid value : 3, 4", tif))
}

func ErrDelistNotAllowed(s string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeDelistNotAllowed, s)
}

func ErrInvalidMarket(s string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidMarket, s)
}

func ErrInvalidCancelTime() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidCancelTime, "Invalid Cancel Time")
}

func ErrNotMatchSender(s string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeNotMatchSender, s)
}

func ErrOrderNotFound(id string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeOrderNotFound, "can not find this order on chain: "+id)
}

func ErrAddressForbidByIssuer() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeAddressForbidByIssuer, "The sender is forbidden by token issuer")
}

func ErrOrderAlreadyExist(id string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeOrderAlreadyExist, "the order [%s] already exist", id)
}

func ErrInvalidOrderAmount(s string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderAmount, s)
}

func ErrNotListedAgainstCet(stock string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeNotListedAgainstCet, "The stock(%s) not have cet trade", stock)
}

func ErrRepeatTradingPair() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeRepeatTradingPair, "The repeatedly created trading pairs")
}

func ErrTokenNoExist() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidToken, "Token not exist")
}

func ErrInvalidOrderID() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderID, "Invalid order id")
}

func ErrInvalidPricePrecision(precision byte) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidPricePrecision, "Invalid price precision : %d", precision)
}

func ErrInvalidPrice(price int64) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidPrice, "Invalid price : %d", price)
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
	return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "Invalid trade pair symbol")
}

func ErrInvalidOrderCommission(err string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderCommission, "The order commission is invalid : %s", err)
}

func ErrStockAndMoneyAreSame() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "Stock and Money should be different")
}

func ErrTokenForbidByIssuer() sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeTokenForbidByIssuer, "Token is frozen by the issuer")
}

func ErrOrderAmountTooSmall(err string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderAmount, "The order amount (%s) too small", err)
}

func ErrDelistRequestExist(market string) sdk.Error {
	return sdk.NewError(CodeSpaceMarket, CodeDelistRequestExist, "The delist request for %s already exists", market)
}
