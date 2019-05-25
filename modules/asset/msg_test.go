package asset

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var tAccAddr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")

func TestMsgIssueToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want string
	}{
		{
			"base-case",
			NewMsgIssueToken("test-coin", "coin", 100000, tAccAddr,
				false, false, false, false),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:           tt.msg.Name,
				Symbol:         tt.msg.Symbol,
				TotalSupply:    tt.msg.TotalSupply,
				Owner:          tt.msg.Owner,
				Mintable:       tt.msg.Mintable,
				Burnable:       tt.msg.Burnable,
				AddrFreezable:  tt.msg.AddrFreezable,
				TokenFreezable: tt.msg.TokenFreezable,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgIssueToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgIssueToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want sdk.Error
	}{
		{
			"case-name1",
			NewMsgIssueToken("123456789012345678901234567890123", "coin", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenName("issue token name limited to 32 unicode characters"),
		},
		{
			"case-symbol1",
			NewMsgIssueToken("name", "1a", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-symbol2",
			NewMsgIssueToken("name", "A999", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-symbol3",
			NewMsgIssueToken("name", "aa1234567", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-symbol4",
			NewMsgIssueToken("name", "a*aa", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-totalSupply1",
			NewMsgIssueToken("name", "coin", 9E18+1, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSupply("issue token supply amt limited to 90 billion"),
		},
		{
			"case-totalSupply2",
			NewMsgIssueToken("name", "coin", -1, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSupply("issue token supply amt should be positive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:           tt.msg.Name,
				Symbol:         tt.msg.Symbol,
				TotalSupply:    tt.msg.TotalSupply,
				Owner:          tt.msg.Owner,
				Mintable:       tt.msg.Mintable,
				Burnable:       tt.msg.Burnable,
				AddrFreezable:  tt.msg.AddrFreezable,
				TokenFreezable: tt.msg.TokenFreezable,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgIssueToken_GetSignBytes(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want string
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, tAccAddr,
				false, false, false, false),
			`{"type":"asset/MsgIssueToken","value":{"addr_freezable":false,"burnable":false,"mintable":false,"name":"ABC Token","owner":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc","token_freezable":false,"total_supply":"100000"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:           tt.msg.Name,
				Symbol:         tt.msg.Symbol,
				TotalSupply:    tt.msg.TotalSupply,
				Owner:          tt.msg.Owner,
				Mintable:       tt.msg.Mintable,
				Burnable:       tt.msg.Burnable,
				AddrFreezable:  tt.msg.AddrFreezable,
				TokenFreezable: tt.msg.TokenFreezable,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgIssueToken.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgIssueToken_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, tAccAddr,
				false, false, false, false),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:           tt.msg.Name,
				Symbol:         tt.msg.Symbol,
				TotalSupply:    tt.msg.TotalSupply,
				Owner:          tt.msg.Owner,
				Mintable:       tt.msg.Mintable,
				Burnable:       tt.msg.Burnable,
				AddrFreezable:  tt.msg.AddrFreezable,
				TokenFreezable: tt.msg.TokenFreezable,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
