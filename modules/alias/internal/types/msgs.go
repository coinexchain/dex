package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

func IsOnlyForCoinEx(alias string) bool {
	if strings.HasPrefix(alias, "coinex") || strings.HasSuffix(alias, "coinex") {
		return true
	}
	if alias == "cet" || alias == "viabtc" || alias == "cetdac" || alias == "www.coinex.com" ||
		alias == "www.coinex.org" {
		return true
	}
	return false
}

func IsValidChar(c rune) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if c == '-' || c == '_' || c == '.' || c == '@' {
		return true
	}
	return false
}

func IsValidAlias(alias string) bool {
	if len(alias) < 2 || len(alias) > 100 {
		return false
	}
	for _, c := range alias {
		if !IsValidChar(c) {
			return false
		}
	}
	return true
}

//=================================

var _ sdk.Msg = MsgAliasUpdate{}

type MsgAliasUpdate struct {
	Owner     sdk.AccAddress `json:"owner"`
	Alias     string         `json:"alias"`
	IsAdd     bool           `json:"is_add"`
	AsDefault bool           `json:"as_default"`
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
