package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the name of the bankx module
const RouterKey = "market"

///////////////////////////////////////////////////////////
// MsgCreateMarketInfo

var _ sdk.Msg = MsgCreateMarketInfo{}

type MsgCreateMarketInfo struct {
	Stock          string `json:"stock"`
	Money          string `json:"money"`
	Crater         string `json:"crater"`
	PricePrecision byte   `json:"priceprecision"`
}

func NewMsgCreateMarketInfo(stock, money, crater string, priceprecision byte) MsgCreateMarketInfo {
	return MsgCreateMarketInfo{Stock: stock, Money: money, Crater: crater, PricePrecision: priceprecision}
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgCreateMarketInfo) Route() string { return RouterKey }

func (msg MsgCreateMarketInfo) Type() string { return "create_market_required" }

func (msg MsgCreateMarketInfo) ValidateBasic() sdk.Error {
	if len(msg.Crater) == 0 {
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
	return []sdk.AccAddress{[]byte(msg.Crater)}
}

///////////////////////////////////////////////////////////
// MsgCreateGTEOrder

var _ sdk.Msg = MsgCreateGTEOrder{}

type MsgCreateGTEOrder struct {
	Sender         string
	Sequence       uint64
	Symbol         string
	OrderType      byte
	PricePrecision byte
	Price          uint64
	Quantity       uint64
	Side           byte
	TimeInForce    int
	Height         uint64
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
	if msg.PricePrecision < 0 || msg.PricePrecision > sdk.Precision {
		return sdk.ErrInvalidAddress("price precision value out of range[0, 18]")
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
