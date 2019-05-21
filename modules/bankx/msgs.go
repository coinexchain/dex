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
