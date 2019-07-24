package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FillDec(s string, dec *sdk.Dec) {
	defer func() {
		if r := recover(); r != nil {
			*dec = sdk.ZeroDec()
		}
	}()
	*dec = sdk.MustNewDecFromStr(s)
}

// /////////////////////////////////////////////////////////

const MaxTradeAmount = int64(10000) * int64(10000) * int64(10000) * int64(10000) * 100 // Ten Billion

var _ sdk.Msg = MsgBancorInit{}
var _ sdk.Msg = MsgBancorTrade{}
var _ sdk.Msg = MsgBancorCancel{}

type MsgBancorInit struct {
	Owner            sdk.AccAddress `json:"owner"`
	Stock            string         `json:"stock"`
	Money            string         `json:"money"`
	MaxSupply        sdk.Int        `json:"max_supply"`
	MaxPrice         sdk.Dec        `json:"max_price"`
	InitPrice        sdk.Dec        `json:"init_price"`
	EnableCancelTime int64          `json:"enable_cancel_time"`
}

type MsgBancorCancel struct {
	Owner sdk.AccAddress `json:"owner"`
	Stock string         `json:"stock"`
	Money string         `json:"money"`
}

type MsgBancorTrade struct {
	Sender     sdk.AccAddress `json:"sender"`
	Stock      string         `json:"stock"`
	Money      string         `json:"money"`
	Amount     int64          `json:"amount"`
	IsBuy      bool           `json:"is_buy"`
	MoneyLimit int64          `json:"money_limit"`
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgBancorInit) Route() string { return RouterKey }

func (msg MsgBancorInit) Type() string { return "bancor_init" }

func (msg MsgBancorInit) ValidateBasic() sdk.Error {
	if len(msg.Owner) == 0 {
		return sdk.ErrInvalidAddress("missing owner address")
	}
	if len(msg.Stock) == 0 || len(msg.Money) == 0 || msg.Stock == "cet" {
		return ErrInvalidSymbol()
	}
	if !msg.MaxSupply.IsPositive() {
		return ErrNonPositiveSupply()
	}
	if !msg.MaxPrice.IsPositive() {
		return ErrNonPositivePrice()
	}
	if msg.InitPrice.IsNegative() {
		return ErrNegativePrice()
	}
	return nil
}

func (msg MsgBancorInit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBancorInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Owner)}
}

func (msg MsgBancorCancel) Route() string { return RouterKey }

func (msg MsgBancorCancel) Type() string { return "bancor_cancel" }

func (msg MsgBancorCancel) ValidateBasic() sdk.Error {
	if len(msg.Owner) == 0 {
		return sdk.ErrInvalidAddress("missing owner address")
	}
	if len(msg.Stock) == 0 || len(msg.Money) == 0 || msg.Stock == "cet" {
		return ErrInvalidSymbol()
	}
	return nil
}

func (msg MsgBancorCancel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBancorCancel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Owner)}
}

func (msg MsgBancorTrade) Route() string { return RouterKey }

func (msg MsgBancorTrade) Type() string { return "bancor_trade" }

func (msg MsgBancorTrade) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.Stock) == 0 || len(msg.Money) == 0 || msg.Stock == "cet" {
		return ErrInvalidSymbol()
	}
	if msg.Amount <= 0 {
		return ErrNonPositiveAmount()
	}
	if msg.Amount > MaxTradeAmount {
		return ErrTradeAmountIsTooLarge()
	}
	return nil
}

func (msg MsgBancorTrade) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBancorTrade) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Sender)}
}
