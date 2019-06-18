package asset

import (
	"os"
	"reflect"
	"testing"

	"github.com/coinexchain/dex/cmd"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMain(m *testing.M) {
	cmd.InitSdkConfig()
	os.Exit(m.Run())
}

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
			ErrorInvalidTokenName("token name is limited to 32 unicode characters"),
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
			ErrorInvalidTokenSupply("token total supply must be positive"),
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
	var addr, _ = sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	tests := []struct {
		name string
		msg  MsgIssueToken
		want string
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, addr,
				false, false, false, false),
			`{"type":"asset/MsgIssueToken","value":{"addr_forbiddable":false,"burnable":false,"mintable":false,"name":"ABC Token","owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","token_forbiddable":false,"total_supply":"100000"}}`,
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
	var addr1, _ = sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	var addr2, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want string
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", addr1, addr2),
			`{"type":"asset/MsgTransferOwnership","value":{"new_owner":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","original_owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
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
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgMintToken
		want string
	}{
		{
			"base-case",
			NewMsgMintToken("abc", 100000, addr),
			`{"type":"asset/MsgMintToken","value":{"amount":"100000","owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc"}}`,
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
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgBurnToken
		want string
	}{
		{
			"base-case",
			NewMsgBurnToken("abc", 100000, addr),
			`{"type":"asset/MsgBurnToken","value":{"amount":"100000","owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc"}}`,
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
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgForbidToken("abc", addr),
			`{"type":"asset/MsgForbidToken","value":{"owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc"}}`,
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
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want string
	}{
		{
			"base-case",
			NewMsgUnForbidToken("abc", addr),
			`{"type":"asset/MsgUnForbidToken","value":{"owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc"}}`,
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

func TestMsgAddTokenWhitelist_Route(t *testing.T) {
	whitelist := mockWhitelist()

	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", tAccAddr, whitelist),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgAddTokenWhitelist.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgAddTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrorInvalidTokenOwner("add token whitelist need a valid owner addr"),
		},
		{
			"case-invalidWhitelist",
			NewMsgAddTokenWhitelist("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidTokenWhitelist("add nil token whitelist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddTokenWhitelist_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	whitelist := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", addr, whitelist),
			`{"type":"asset/MsgAddTokenWhitelist","value":{"owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgAddTokenWhitelist.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgAddTokenWhitelist_GetSigners(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", tAccAddr, whitelist),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddTokenWhitelist.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_Route(t *testing.T) {
	whitelist := mockWhitelist()

	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", tAccAddr, whitelist),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgRemoveTokenWhitelist.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgRemoveTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrorInvalidTokenOwner("remove token whitelist need a valid owner addr"),
		},
		{
			"case-invalidWhitelist",
			NewMsgRemoveTokenWhitelist("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidTokenWhitelist("remove nil token whitelist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	whitelist := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want string
	}{
		{
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", addr, whitelist),
			`{"type":"asset/MsgRemoveTokenWhitelist","value":{"owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_GetSigners(t *testing.T) {
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", tAccAddr, whitelist),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveTokenWhitelist{
				tt.msg.Symbol,
				tt.msg.OwnerAddress,
				tt.msg.Whitelist,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_Route(t *testing.T) {
	addr := mockAddresses()

	tests := []struct {
		name string
		msg  MsgForbidAddr
		want string
	}{
		{
			"base-case",
			NewMsgForbidAddr("abc", tAccAddr, addr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.ForbidAddr,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgForbidAddr.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_ValidateBasic(t *testing.T) {
	addr := mockAddresses()
	tests := []struct {
		name string
		msg  MsgForbidAddr
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgForbidAddr("abc", sdk.AccAddress{}, addr),
			ErrorInvalidTokenOwner("forbid address need a valid owner addr"),
		},
		{
			"case-invalidAddr",
			NewMsgForbidAddr("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidAddress("forbid nil address"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.ForbidAddr,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	addresses := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgForbidAddr
		want string
	}{
		{
			"base-case",
			NewMsgForbidAddr("abc", addr, addresses),
			`{"type":"asset/MsgForbidAddr","value":{"forbid_addr":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"],"owner_address":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.ForbidAddr,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgForbidAddr.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_GetSigners(t *testing.T) {
	addr := mockAddresses()
	tests := []struct {
		name string
		msg  MsgForbidAddr
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgForbidAddr("abc", tAccAddr, addr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.ForbidAddr,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidAddr.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_Route(t *testing.T) {
	addr := mockAddresses()

	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want string
	}{
		{
			"base-case",
			NewMsgUnForbidAddr("abc", tAccAddr, addr),
			RouterKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.UnForbidAddr,
			}
			if got := msg.Route(); got != tt.want {
				t.Errorf("MsgUnForbidAddr.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_ValidateBasic(t *testing.T) {
	addr := mockAddresses()
	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want sdk.Error
	}{
		{
			"case-invalidOwner",
			NewMsgUnForbidAddr("abc", sdk.AccAddress{}, addr),
			ErrorInvalidTokenOwner("unforbid address need a valid owner addr"),
		},
		{
			"case-invalidAddr",
			NewMsgUnForbidAddr("abc", tAccAddr, []sdk.AccAddress{}),
			ErrorInvalidAddress("unforbid nil address"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.UnForbidAddr,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_GetSignBytes(t *testing.T) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	addresses := []sdk.AccAddress{addr1, addr2}
	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want string
	}{
		{
			"base-case",
			NewMsgUnForbidAddr("abc", addr, addresses),
			`{"type":"asset/MsgUnForbidAddr","value":{"owner_addr":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","symbol":"abc","unforbid_addr":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.UnForbidAddr,
			}
			if got := msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MsgUnForbidAddr.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_GetSigners(t *testing.T) {
	addr := mockAddresses()
	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want []sdk.AccAddress
	}{
		{
			"base-case",
			NewMsgUnForbidAddr("abc", tAccAddr, addr),
			[]sdk.AccAddress{tAccAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgUnForbidAddr{
				tt.msg.Symbol,
				tt.msg.OwnerAddr,
				tt.msg.UnForbidAddr,
			}
			if got := msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidAddr.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}
