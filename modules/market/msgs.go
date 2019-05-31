package market

import (
	"strconv"
	"strings"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the name of the bankx module
const (
	RouterKey = "market"
	MarketKey = RouterKey
)

var (
	msgCdc = codec.New()
)

func init() {
	RegisterCodec(msgCdc)
}

///////////////////////////////////////////////////////////
// MsgCreateMarketInfo

var _ sdk.Msg = MsgCreateMarketInfo{}

type MsgCreateMarketInfo struct {
	Stock          string         `json:"stock"`
	Money          string         `json:"money"`
	Creator        sdk.AccAddress `json:"creator"`
	PricePrecision byte           `json:"price_precision"`
}

func NewMsgCreateMarketInfo(stock, money string, crater sdk.AccAddress, priceprecision byte) MsgCreateMarketInfo {
	return MsgCreateMarketInfo{Stock: stock, Money: money, Creator: crater, PricePrecision: priceprecision}
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgCreateMarketInfo) Route() string { return RouterKey }

func (msg MsgCreateMarketInfo) Type() string { return "create_market_required" }

func (msg MsgCreateMarketInfo) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("missing creator address")
	}
	if len(msg.Stock) == 0 || len(msg.Money) == 0 {
		return sdk.ErrInvalidAddress("missing stock or money identifier")
	}
	if msg.PricePrecision > sdk.Precision {
		return sdk.ErrInvalidAddress("proceprecision value out of range[0, 18]")
	}
	return nil
}

func (msg MsgCreateMarketInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateMarketInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Creator)}
}

///////////////////////////////////////////////////////////
// MsgCreateOrder

var _ sdk.Msg = MsgCreateOrder{}

type MsgCreateOrder struct {
	Sender         sdk.AccAddress `json:"sender"`
	Sequence       uint64         `json:"sequence"`
	Symbol         string         `json:"symbol"`
	OrderType      byte           `json:"ordertype"`
	PricePrecision byte           `json:"priceprecision"`
	Price          int64          `json:"price"`
	Quantity       int64          `json:"quantity"`
	Side           byte           `json:"side"`
	TimeInForce    int            `json:"timeinforce"`
}

func (msg MsgCreateOrder) Route() string { return RouterKey }

func (msg MsgCreateOrder) Type() string { return "create_order_required" }

func (msg MsgCreateOrder) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("missing creator address")
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrInvalidAddress("missing GTE order symbol identifier")
	}
	if msg.PricePrecision < MinimumTokenPricePrecision ||
		msg.PricePrecision > MaxTokenPricePrecision {
		return sdk.ErrInvalidAddress("price precision value out of range[8, 18]. actual : " + strconv.Itoa(int(msg.PricePrecision)))
	}

	if msg.Side != match.BUY && msg.Side != match.SELL {
		return ErrInvalidTradeSide()
	}

	if msg.OrderType != LimitOrder {
		return ErrInvalidOrderType()
	}

	if len(strings.Split(msg.Symbol, SymbolSeparator)) != 2 {
		return ErrInvalidSymbol()
	}

	if msg.Price < 0 || msg.Price > asset.MaxTokenAmount {
		return ErrInvalidPrice()
	}

	//TODO: Add Other checks
	return nil
}

func (msg MsgCreateOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Sender)}
}

func (msg MsgCreateOrder) IsGTEOrder() bool {
	return msg.TimeInForce == GTE
}

///////////////////////////////////////////////////////////
// MsgCancelOrder

type MsgCancelOrder struct {
	Sender  sdk.AccAddress `json:"sender"`
	OrderID string         `json:"order_id"`
}

func (msg MsgCancelOrder) Route() string {
	return MarketKey
}

func (msg MsgCancelOrder) Type() string {
	return "cancel_order_required"
}

func (msg MsgCancelOrder) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress()
	}

	if len(strings.Split(msg.OrderID, "-")) != 2 {
		return ErrInvalidOrderID()
	}

	return nil
}

func (msg MsgCancelOrder) GetSignBytes() []byte {

	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelOrder) GetSigners() []sdk.AccAddress {

	return []sdk.AccAddress{msg.Sender}
}
