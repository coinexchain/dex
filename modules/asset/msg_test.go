package asset

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgIssueToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want string
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC token", "abc", 100000, tAccAddr,
				false, false, false, false),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:             tt.msg.Name,
				Symbol:           tt.msg.Symbol,
				TotalSupply:      tt.msg.TotalSupply,
				Owner:            tt.msg.Owner,
				Mintable:         tt.msg.Mintable,
				Burnable:         tt.msg.Burnable,
				AddrForbiddable:  tt.msg.AddrForbiddable,
				TokenForbiddable: tt.msg.TokenForbiddable,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgIssueToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgIssueToken_ValidateBasic(t *testing.T) {
	invalidTokenSymbol := ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}")
	tests := []struct {
		name string
		msg  MsgIssueToken
		want sdk.Error
	}{
		{
			"case-name1",
			NewMsgIssueToken("123456789012345678901234567890123", "coin", 100000, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenName("token name limited to 32 unicode characters"),
		},
		{
			"case-symbol1",
			NewMsgIssueToken("name", "1a", 100000, tAccAddr,
				false, false, false, false),
			invalidTokenSymbol,
		},
		{
			"case-symbol2",
			NewMsgIssueToken("name", "A999", 100000, tAccAddr,
				false, false, false, false),
			invalidTokenSymbol,
		},
		{
			"case-symbol3",
			NewMsgIssueToken("name", "aa1234567", 100000, tAccAddr,
				false, false, false, false),
			invalidTokenSymbol,
		},
		{
			"case-symbol4",
			NewMsgIssueToken("name", "a*aa", 100000, tAccAddr,
				false, false, false, false),
			invalidTokenSymbol,
		},
		{
			"case-totalSupply1",
			NewMsgIssueToken("name", "coin", 9E18+1, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSupply("token total supply before 1e8 boosting should be less than 90 billion"),
		},
		{
			"case-totalSupply2",
			NewMsgIssueToken("name", "coin", -1, tAccAddr,
				false, false, false, false),
			ErrorInvalidTokenSupply("token total supply must a positive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:             tt.msg.Name,
				Symbol:           tt.msg.Symbol,
				TotalSupply:      tt.msg.TotalSupply,
				Owner:            tt.msg.Owner,
				Mintable:         tt.msg.Mintable,
				Burnable:         tt.msg.Burnable,
				AddrForbiddable:  tt.msg.AddrForbiddable,
				TokenForbiddable: tt.msg.TokenForbiddable,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgIssueToken_GetSignBytes(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	tests := []struct {
		name string
		msg  MsgIssueToken
		want string
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, addr,
				false, false, false, false),
			`{"type":"asset/MsgIssueToken","value":{"addr_forbiddable":false,"burnable":false,"mintable":false,"name":"ABC Token","owner":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc","token_forbiddable":false,"total_supply":"100000"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgIssueToken{
				Name:             tt.msg.Name,
				Symbol:           tt.msg.Symbol,
				TotalSupply:      tt.msg.TotalSupply,
				Owner:            tt.msg.Owner,
				Mintable:         tt.msg.Mintable,
				Burnable:         tt.msg.Burnable,
				AddrForbiddable:  tt.msg.AddrForbiddable,
				TokenForbiddable: tt.msg.TokenForbiddable,
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
				Name:             tt.msg.Name,
				Symbol:           tt.msg.Symbol,
				TotalSupply:      tt.msg.TotalSupply,
				Owner:            tt.msg.Owner,
				Mintable:         tt.msg.Mintable,
				Burnable:         tt.msg.Burnable,
				AddrForbiddable:  tt.msg.AddrForbiddable,
				TokenForbiddable: tt.msg.TokenForbiddable,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want string
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", tAccAddr, tAccAddr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgTransferOwnership{
				tt.msg.Symbol,
				tt.msg.OriginalOwner,
				tt.msg.NewOwner,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgTransferOwnership.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want sdk.Error
	}{
		{
			"case-invalid1",
			NewMsgTransferOwnership("abc", sdk.AccAddress{}, tAccAddr),
			ErrorInvalidTokenOwner("transfer owner ship need a valid addr"),
		},
		{
			"case-invalid2",
			NewMsgTransferOwnership("abc", tAccAddr, sdk.AccAddress{}),
			ErrorInvalidTokenOwner("transfer owner ship need a valid addr"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgTransferOwnership{
				tt.msg.Symbol,
				tt.msg.OriginalOwner,
				tt.msg.NewOwner,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgTransferOwnership.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	var addr2, _ = sdk.AccAddressFromBech32("cosmos1r8rjvkawsq379z7qndtqtkks0pvqxxepnk0frr")
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want string
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", addr1, addr2),
			`{"type":"asset/MsgTransferOwnership","value":{"new_owner":"cosmos1r8rjvkawsq379z7qndtqtkks0pvqxxepnk0frr","original_owner":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgTransferOwnership{
				tt.msg.Symbol,
				tt.msg.OriginalOwner,
				tt.msg.NewOwner,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgTransferOwnership.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", tAccAddr, sdk.AccAddress{}),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgTransferOwnership{
				tt.msg.Symbol,
				tt.msg.OriginalOwner,
				tt.msg.NewOwner,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgTransferOwnership.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgMintToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMintToken
		want string
	}{
		{
			"base-case",
			NewMsgMintToken("abc", 1000000, tAccAddr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgMintToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgMintToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgMintToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMintToken
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgMintToken("abc", 10000, sdk.AccAddress{}),
			ErrorInvalidTokenOwner("mint token need a valid owner addr"),
		},
		{
			"case-invalidAmt1",
			NewMsgMintToken("abc", 9E18+1, tAccAddr),
			ErrorInvalidTokenMint("token total supply before 1e8 boosting should be less than 90 billion"),
		},
		{
			"case-invalidAmt2",
			NewMsgMintToken("abc", -1, tAccAddr),
			ErrorInvalidTokenMint("mint amount should be positive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgMintToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgMintToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgMintToken_GetSignBytes(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	tests := []struct {
		name string
		msg  MsgMintToken
		want string
	}{
		{
			"base-case",
			NewMsgMintToken("abc", 100000, addr),
			`{"type":"asset/MsgMintToken","value":{"amount":"100000","owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgMintToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgMintToken.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgMintToken_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMintToken
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgMintToken("abc", 100000, tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgMintToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgMintToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBurnToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBurnToken
		want string
	}{
		{
			"base-case",
			NewMsgBurnToken("abc", 1000000, tAccAddr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBurnToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgBurnToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBurnToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBurnToken
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgBurnToken("abc", 10000, sdk.AccAddress{}),
			ErrorInvalidTokenOwner("burn token need a valid owner addr"),
		},
		{
			"case-invalidAmt1",
			NewMsgBurnToken("abc", 9E18+1, tAccAddr),
			ErrorInvalidTokenBurn("token total supply before 1e8 boosting should be less than 90 billion"),
		},
		{
			"case-invalidAmt2",
			NewMsgBurnToken("abc", -1, tAccAddr),
			ErrorInvalidTokenBurn("burn amount should be positive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBurnToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBurnToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBurnToken_GetSignBytes(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	tests := []struct {
		name string
		msg  MsgBurnToken
		want string
	}{
		{
			"base-case",
			NewMsgBurnToken("abc", 100000, addr),
			`{"type":"asset/MsgBurnToken","value":{"amount":"100000","owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBurnToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgBurnToken.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgBurnToken_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBurnToken
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgBurnToken("abc", 100000, tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBurnToken{
				tt.msg.Symbol,
				tt.msg.Amount,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBurnToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgForbidToken("abc", tAccAddr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgForbidToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgForbidToken
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgForbidToken("abc", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("forbid token need a valid owner addr"),
		},
		{
			"case-invalidSymbol",
			NewMsgForbidToken("*90", sdk.AccAddress{}),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidToken_GetSignBytes(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	tests := []struct {
		name string
		msg  MsgForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgForbidToken("abc", addr),
			`{"type":"asset/MsgForbidToken","value":{"owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgForbidToken.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgForbidToken_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgForbidToken
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgForbidToken("abc", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidToken_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgUnForbidToken("abc", tAccAddr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgUnForbidToken.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgUnForbidToken("abc", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("forbid token need a valid owner addr"),
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidToken("*90", sdk.AccAddress{}),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidToken_GetSignBytes(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgUnForbidToken("abc", addr),
			`{"type":"asset/MsgUnForbidToken","value":{"owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgUnForbidToken.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidToken_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgUnForbidToken("abc", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidToken{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidToken.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddForbidWhitelist_Route(t *testing.T) {
	whitelist := mockWhitelist()

	tests := []struct {
		name string
		msg  MsgAddForbidWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgAddForbidWhitelist("abc", tAccAddr, whitelist),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgAddForbidWhitelist.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddForbidWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgAddForbidWhitelist
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgAddForbidWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrorInvalidTokenOwner("add forbid whitelist need a valid owner addr"),
		},
		{
			"case-invalidWhitelist",
			NewMsgAddForbidWhitelist("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidTokenWhitelist("add nil forbid whitelist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddForbidWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddForbidWhitelist_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt")
	var addr2, _ = sdk.AccAddressFromBech32("cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf")
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	whitelist := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgAddForbidWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgAddForbidWhitelist("abc", addr, whitelist),
			`{"type":"asset/MsgAddForbidWhitelist","value":{"owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc","whitelist":["cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt","cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgAddForbidWhitelist.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgAddForbidWhitelist_GetSigners(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgAddForbidWhitelist
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgAddForbidWhitelist("abc", tAccAddr, whitelist),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddForbidWhitelist.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveForbidWhitelist_Route(t *testing.T) {
	whitelist := mockWhitelist()

	tests := []struct {
		name string
		msg  MsgRemoveForbidWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgRemoveForbidWhitelist("abc", tAccAddr, whitelist),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgRemoveForbidWhitelist.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveForbidWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgRemoveForbidWhitelist
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgRemoveForbidWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrorInvalidTokenOwner("remove forbid whitelist need a valid owner addr"),
		},
		{
			"case-invalidWhitelist",
			NewMsgRemoveForbidWhitelist("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidTokenWhitelist("remove nil forbid whitelist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveForbidWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveForbidWhitelist_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt")
	var addr2, _ = sdk.AccAddressFromBech32("cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf")
	var addr, _ = sdk.AccAddressFromBech32("cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd")
	whitelist := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgRemoveForbidWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgRemoveForbidWhitelist("abc", addr, whitelist),
			`{"type":"asset/MsgRemoveForbidWhitelist","value":{"owner_address":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","symbol":"abc","whitelist":["cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt","cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgRemoveForbidWhitelist.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveForbidWhitelist_GetSigners(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgRemoveForbidWhitelist
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgRemoveForbidWhitelist("abc", tAccAddr, whitelist),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveForbidWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveForbidWhitelist.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
