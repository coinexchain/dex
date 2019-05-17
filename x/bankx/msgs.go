package bankx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the name of the bankx module
const RouterKey = "bankx"

var _ sdk.Msg = MsgSetTransferMemoRequired{}

type MsgSetTransferMemoRequired struct {
	Address  sdk.AccAddress `json:"address"`
	Required bool           `json:"required"`
}

func NewMsgSetTransferMemoRequired(addr sdk.AccAddress, required bool) MsgSetTransferMemoRequired {
	return MsgSetTransferMemoRequired{Address: addr, Required: required}
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgSetTransferMemoRequired) Route() string { return RouterKey }

func (msg MsgSetTransferMemoRequired) Type() string { return "set_transfer_memo_required" }

func (msg MsgSetTransferMemoRequired) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress("missing address")
	}
	return nil
}

func (msg MsgSetTransferMemoRequired) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetTransferMemoRequired) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
