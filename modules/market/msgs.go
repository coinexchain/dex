package market

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/match"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// RouterKey is the name of the bankx module
const (
	RouterKey = "market"
	MarketKey = RouterKey
)

///////////////////////////////////////////////////////////
// MsgCreateMarketInfo

var _ sdk.Msg = MsgCreateMarketInfo{}

type MsgCreateMarketInfo struct {
	Stock          string         `json:"stock"`
	Money          string         `json:"money"`
	Creator        sdk.AccAddress `json:"crater"`
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
	if msg.PricePrecision < 0 || msg.PricePrecision > sdk.Precision {
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
// MsgCreateGTEOrder

var _ sdk.Msg = MsgCreateGTEOrder{}

type MsgCreateGTEOrder struct {
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

func (msg MsgCreateGTEOrder) Route() string { return RouterKey }

func (msg MsgCreateGTEOrder) Type() string { return "create_gte_order_required" }

func (msg MsgCreateGTEOrder) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("missing creator address")
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrInvalidAddress("missing GTE order symbol identifier")
	}
	if msg.PricePrecision < MinimumTokenPricePrecision ||
		msg.PricePrecision > MaxTokenPricePrecision {
		return sdk.ErrInvalidAddress("price precision value out of range[0, 18]")
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

func (msg MsgCreateGTEOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateGTEOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Sender)}
}

///////////////////////////////////////////////////////////
// MsgCreateIOCOrder

type MsgCreateIOCOrder struct {
}

///////////////////////////////////////////////////////////
// MsgCancelOrder

type MsgCancelOrder struct {
}
