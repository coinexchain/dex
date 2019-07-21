package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"unicode"
	"unicode/utf8"
)

const AliasMaxLength = 20

func IsValidAlias(alias string) bool {
	if len(alias) == 0 {
		return false
	}
	if !utf8.ValidString(alias) {
		return false
	}
	for i, c := range alias {
		if i == 0 && '0' <= c && c <= '9' {
			return false
		}
		if unicode.IsSpace(c) || unicode.IsControl(c) {
			return false
		}
		if i >= AliasMaxLength {
			return false
		}
	}
	return true
}

//=================================

var _ sdk.Msg = MsgAliasUpdate{}

type MsgAliasUpdate struct {
	Owner sdk.AccAddress `json:"owner"`
	Alias string         `json:"alias"`
	IsAdd bool           `json:"is_add"`
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
	return []sdk.AccAddress{[]byte(msg.Owner)}
}
