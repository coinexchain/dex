package bankx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
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

var _ sdk.Msg = MsgSendWithUnlockTime{}

type MsgSendWithUnlockTime struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
	UnlockTime  int64          `json:"unlock_time"`
}

func NewMsgSendWithUnlocktime(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins, time int64) MsgSendWithUnlockTime {
	return MsgSendWithUnlockTime{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount, UnlockTime: time}
}

func (msg MsgSendWithUnlockTime) Route() string {
	return RouterKey
}

func (msg MsgSendWithUnlockTime) Type() string {
	return "unlock_time_send"
}

func (msg MsgSendWithUnlockTime) ValidateBasic() sdk.Error {
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
	t := time.Now().Unix()
	if msg.UnlockTime <= t {
		return ErrUnlockTime("Invalid Unlock Time")
	}
	return nil
}

func (msg MsgSendWithUnlockTime) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSendWithUnlockTime) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
