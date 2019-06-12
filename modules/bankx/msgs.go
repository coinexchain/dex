package bankx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the name of the bankx module
const RouterKey = "bankx"

var _ sdk.Msg = MsgSetMemoRequired{}

type MsgSetMemoRequired struct {
	Address  sdk.AccAddress `json:"address"`
	Required bool           `json:"required"`
}

func NewMsgSetTransferMemoRequired(addr sdk.AccAddress, required bool) MsgSetMemoRequired {
	return MsgSetMemoRequired{Address: addr, Required: required}
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgSetMemoRequired) Route() string { return RouterKey }

func (msg MsgSetMemoRequired) Type() string { return "set_memo_required" }

func (msg MsgSetMemoRequired) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress("missing address")
	}
	return nil
}

func (msg MsgSetMemoRequired) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetMemoRequired) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

var _ sdk.Msg = MsgSend{}

type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
	UnlockTime  int64          `json:"unlock_time"`
}

func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins, time int64) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount, UnlockTime: time}
}

func (msg MsgSend) Route() string {
	return RouterKey
}

func (msg MsgSend) Type() string {
	return "send"
}

func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	if msg.UnlockTime < 0 {
		return ErrUnlockTime("negative unlock time ")
	}

	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
