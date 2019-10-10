package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgAliasUpdate{}

type MsgAliasUpdate struct {
	Owner     sdk.AccAddress `json:"owner"`
	Alias     string         `json:"alias"`
	IsAdd     bool           `json:"is_add"`
	AsDefault bool           `json:"as_default"`
}

func (msg *MsgAliasUpdate) SetAccAddress(addr sdk.AccAddress) {
	msg.Owner = addr
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgAliasUpdate) Route() string { return RouterKey }

func (msg MsgAliasUpdate) Type() string { return "alias_update" }

func (msg MsgAliasUpdate) ValidateBasic() sdk.Error {
	if len(msg.Owner) == 0 {
		return sdk.ErrInvalidAddress("missing owner address")
	}
	if len(msg.Alias) == 0 {
		return ErrEmptyAlias()
	}
	if !IsValidAlias(msg.Alias) {
		return ErrInvalidAlias()
	}
	return nil
}

func (msg MsgAliasUpdate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAliasUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
