package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the name of the market module
const (
	// msg keys for Kafka
	CreateTradingInfoKey = "create_trading_info"
	CancelTradingInfoKey = "cancel_trading_info"

	CreateOrderInfoKey    = "create_order_info"
	FillOrderInfoKey      = "fill_order_info"
	CancelOrderInfoKey    = "del_order_info"
	PricePrecisionInfoKey = "modify-price-precision"
)

// cancel order of reasons
const (
	CancelOrderByManual        = "Manually cancel the order"
	CancelOrderByAllFilled     = "The order was fully filled"
	CancelOrderByGteTimeOut    = "GTE order timeout"
	CancelOrderByIocType       = "IOC order cancel "
	CancelOrderByNoEnoughMoney = "Insufficient freeze money"
	CancelOrderByNotKnow       = "Don't know"
)

// /////////////////////////////////////////////////////////
// MsgCreateTradingPair

var _ sdk.Msg = MsgCreateTradingPair{}

type MsgCreateTradingPair struct {
	Stock          string         `json:"stock"`
	Money          string         `json:"money"`
	Creator        sdk.AccAddress `json:"creator"`
	PricePrecision byte           `json:"price_precision"`
}

func NewMsgCreateTradingPair(stock, money string, crater sdk.AccAddress, pricePrecision byte) MsgCreateTradingPair {
	return MsgCreateTradingPair{
		Stock:          stock,
		Money:          money,
		Creator:        crater,
		PricePrecision: pricePrecision,
	}
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgCreateTradingPair) SetAccAddress(address sdk.AccAddress) {
	msg.Creator = address
}

func (msg MsgCreateTradingPair) Route() string { return RouterKey }

func (msg MsgCreateTradingPair) Type() string { return "create_market_info" }

func (msg MsgCreateTradingPair) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("missing creator address")
	}
	if len(msg.Stock) == 0 || len(msg.Money) == 0 {
		return ErrInvalidSymbol()
	}
	if msg.PricePrecision < MinTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision()
	}
	if msg.Money == msg.Stock {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "stock and money should be different")
	}
	return nil
}

func (msg MsgCreateTradingPair) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateTradingPair) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Creator)}
}

// /////////////////////////////////////////////////////////
// MsgCreateOrder

var _ sdk.Msg = MsgCreateOrder{}

type MsgCreateOrder struct {
	Sender         sdk.AccAddress `json:"sender"`
	Identify       byte           `json:"identify"`
	TradingPair    string         `json:"trading_pair"`
	OrderType      byte           `json:"order_type"`
	PricePrecision byte           `json:"price_precision"`
	Price          int64          `json:"price"`
	Quantity       int64          `json:"quantity"`
	Side           byte           `json:"side"`
	TimeInForce    int            `json:"time_in_force"`
	ExistBlocks    int            `json:"exist_blocks"`
}

func (msg MsgCreateOrder) SetAccAddress(address sdk.AccAddress) {
	msg.Sender = address
}

func (msg MsgCreateOrder) Route() string { return RouterKey }

func (msg MsgCreateOrder) Type() string { return "create_order" }

func (msg MsgCreateOrder) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("missing creator address")
	}
	if len(msg.TradingPair) == 0 {
		return sdk.ErrInvalidAddress("missing GTE order TradingPair identifier")
	}
	if len(strings.Split(msg.TradingPair, SymbolSeparator)) != 2 {
		return ErrInvalidSymbol()
	}
	if msg.OrderType != LimitOrder {
		return ErrInvalidOrderType()
	}
	if msg.PricePrecision < MinTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return sdk.ErrInvalidAddress(fmt.Sprintf("price precision value out of range [0, 18]. actual : %d", msg.PricePrecision))
	}
	if msg.Price <= 0 {
		return ErrInvalidPrice(msg.Price)
	}
	if msg.Quantity < 0 {
		return ErrOrderQuantityToSmall(fmt.Sprintf("%d", msg.Quantity))
	}
	if msg.Side != BUY && msg.Side != SELL {
		return ErrInvalidTradeSide()
	}
	if msg.TimeInForce != GTE && msg.TimeInForce != IOC {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidTimeInforce, fmt.Sprintf("Invalid timeInforce : %d; The valid value : 3, 4", msg.TimeInForce))
	}
	if msg.ExistBlocks < 0 {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidExistBlocks, fmt.Sprintf("Invalid existence time : %d; The range of expected values [0, +∞] ", msg.ExistBlocks))
	}
	if msg.Identify < 0 || msg.Identify > 255 {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidOrderExist, fmt.Sprintf("invalid identify : %d, expected range [0, 255]", msg.Identify))
	}

	return nil
}

func (msg MsgCreateOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Sender)}
}

func (msg MsgCreateOrder) IsGTEOrder() bool {
	return msg.TimeInForce == GTE
}

// /////////////////////////////////////////////////////////
// MsgCancelOrder

type MsgCancelOrder struct {
	Sender  sdk.AccAddress `json:"sender"`
	OrderID string         `json:"order_id"`
}

func (msg MsgCancelOrder) SetAccAddress(address sdk.AccAddress) {
	msg.Sender = address
}

func (msg MsgCancelOrder) Route() string {
	return StoreKey
}

func (msg MsgCancelOrder) Type() string {
	return "cancel_order"
}

func (msg MsgCancelOrder) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress()
	}
	if len(strings.Split(msg.OrderID, OrderIDSeparator)) != OrderIDPartsNum {
		return ErrInvalidOrderID()
	}

	return nil
}

func (msg MsgCancelOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// /////////////////////////////////////////////////////////
// MsgCancelTradingPair

type MsgCancelTradingPair struct {
	Sender        sdk.AccAddress `json:"sender"`
	TradingPair   string         `json:"trading_pair"`
	EffectiveTime int64          `json:"effective_height"`
}

func (msg MsgCancelTradingPair) SetAccAddress(address sdk.AccAddress) {
	msg.Sender = address
}

func (msg MsgCancelTradingPair) Route() string {
	return RouterKey
}

func (msg MsgCancelTradingPair) Type() string {
	return "cancel_market"
}

func (msg MsgCancelTradingPair) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress()
	}
	if len(strings.Split(msg.TradingPair, SymbolSeparator)) != 2 {
		return ErrInvalidSymbol()
	}
	if msg.EffectiveTime < 0 {
		return sdk.NewError(CodeSpaceMarket, CodeInvalidTime, "Invalid height")
	}

	return nil
}

func (msg MsgCancelTradingPair) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelTradingPair) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// -------------------------------------------------
// MsgModifyPricePrecision

type MsgModifyPricePrecision struct {
	Sender         sdk.AccAddress `json:"sender"`
	TradingPair    string         `json:"trading_pair"`
	PricePrecision byte           `json:"price_precision"`
}

func (msg MsgModifyPricePrecision) SetAccAddress(address sdk.AccAddress) {
	msg.Sender = address
}

func (msg MsgModifyPricePrecision) Route() string {
	return RouterKey
}

func (msg MsgModifyPricePrecision) Type() string {
	return "modify_trading_pair_price_precision"
}

func (msg MsgModifyPricePrecision) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return ErrInvalidAddress()
	}
	if len(strings.Split(msg.TradingPair, SymbolSeparator)) != 2 {
		return ErrInvalidSymbol()
	}
	if msg.PricePrecision < MinTokenPricePrecision || msg.PricePrecision > MaxTokenPricePrecision {
		return ErrInvalidPricePrecision()
	}

	return nil
}

func (msg MsgModifyPricePrecision) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgModifyPricePrecision) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// --------------------------------------------------------------------
// msg queue infos for kafka

type CreateMarketInfo struct {
	Stock          string `json:"stock"`
	Money          string `json:"money"`
	PricePrecision byte   `json:"price_precision"`

	// create market info
	Creator      string `json:"creator"`
	CreateHeight int64  `json:"create_height"`
}

type CancelMarketInfo struct {
	Stock string `json:"stock"`
	Money string `json:"money"`

	// del market info
	Deleter string `json:"deleter"`
	DelTime int64  `json:"del_time"`
}

type CreateOrderInfo struct {
	OrderID     string  `json:"order_id"`
	Sender      string  `json:"sender"`
	TradingPair string  `json:"trading_pair"`
	OrderType   byte    `json:"order_type"`
	Price       sdk.Dec `json:"price"`
	Quantity    int64   `json:"quantity"`
	Side        byte    `json:"side"`
	TimeInForce int     `json:"time_in_force"`
	FeatureFee  int64   `json:"feature_fee"`
	Height      int64   `json:"height"`
	FrozenFee   int64   `json:"frozen_fee"`
	Freeze      int64   `json:"freeze"`
}

type FillOrderInfo struct {
	OrderID     string  `json:"order_id"`
	TradingPair string  `json:"trading_pair"`
	Height      int64   `json:"height"`
	Side        byte    `json:"side"`
	Price       sdk.Dec `json:"price"`

	// These fields will change when order was filled/canceled.
	LeftStock int64 `json:"left_stock"`
	Freeze    int64 `json:"freeze"`
	DealStock int64 `json:"deal_stock"`
	DealMoney int64 `json:"deal_money"`
	CurrStock int64 `json:"curr_stock"`
	CurrMoney int64 `json:"curr_money"`
}

type CancelOrderInfo struct {
	OrderID     string  `json:"order_id"`
	TradingPair string  `json:"trading_pair"`
	Height      int64   `json:"height"`
	Side        byte    `json:"side"`
	Price       sdk.Dec `json:"price"`

	// Del infos
	DelReason string `json:"del_reason"`

	// Fields of amount
	UsedCommission int64 `json:"used_commission"`
	LeftStock      int64 `json:"left_stock"`
	RemainAmount   int64 `json:"remain_amount"`
	DealStock      int64 `json:"deal_stock"`
	DealMoney      int64 `json:"deal_money"`
}

type ModifyPricePrecisionInfo struct {
	Sender            string `json:"sender"`
	TradingPair       string `json:"trading_pair"`
	OldPricePrecision byte   `json:"old_price_precision"`
	NewPricePrecision byte   `json:"new_price_precision"`
}
