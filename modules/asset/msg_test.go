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
			"test1",
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

			switch tt.name {
			case "test1":
				if got := msg.Route(); got != tt.want {
					t.Errorf("MsgIssueToken.Route() = %v, want %v", got, tt.want)
				}

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
			"testName1",
			NewMsgIssueToken("123456789012345678901234567890123", "coin", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenName("issue token name limited to 32 unicode characters"),
		},
		{
			"testSymbol1",
			NewMsgIssueToken("name", "1a", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"testSymbol2",
			NewMsgIssueToken("name", "A999", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"testSymbol3",
			NewMsgIssueToken("name", "aa1234567", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"testSymbol4",
			NewMsgIssueToken("name", "a*aa", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSymbol("issue token symbol limited to [a-z][a-z0-9]{1,7}"),
		},
		{
			"testTotalSupply1",
			NewMsgIssueToken("name", "coin", 9E18+1, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSupply("issue token supply amt limited to 90 billion"),
		},
		{
			"testTotalSupply2",
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
			switch tt.name {
			case "testName1":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testSymbol1":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testSymbol2":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testSymbol3":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testSymbol4":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testTotalSupply1":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
			case "testTotalSupply2":
				if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
				}
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
			"test1",
			NewMsgIssueToken("test-coin", "coin", 100000, tAccAddr,
				false, false, false, false),
			`{"type":"asset/MsgIssueToken","value":{"AddrFreezable":false,"Burnable":false,"Mintable":false,"Name":"test-coin","Owner":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","Symbol":"coin","TokenFreezable":false,"TotalSupply":"100000"}}`,
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
			"test1",
			NewMsgIssueToken("test-coin", "coin", 100000, tAccAddr,
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
