package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgIssueToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIssueToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			nil,
		},
		{
			"case-name",
			NewMsgIssueToken(string(make([]byte, 32+1)), "abc", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenName(string(make([]byte, 32+1))),
		},
		{
			"case-owner",
			NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(10000), sdk.AccAddress{},
				false, false, false, false, "", "", TestIdentityString),
			ErrNilTokenOwner(),
		},
		{
			"case-symbol1",
			NewMsgIssueToken("ABC Token", "1aa", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenSymbol("1aa"),
		},
		{
			"case-symbol2",
			NewMsgIssueToken("ABC Token", "A999", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenSymbol("A999"),
		},
		{
			"case-symbol3",
			NewMsgIssueToken("ABC Token", "aa345678901234567", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenSymbol("aa345678901234567"),
		},
		{
			"case-symbol4",
			NewMsgIssueToken("ABC Token", "a*aa", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenSymbol("a*aa"),
		},
		{
			"case-totalSupply",
			NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(-1), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			ErrInvalidTokenSupply(sdk.NewInt(-1).String()),
		},
		{
			"case-url",
			NewMsgIssueToken("name", "coin", sdk.NewInt(10000), testAddr,
				false, false, false, false, string(make([]byte, MaxTokenURLLength+1)), "", TestIdentityString),
			ErrInvalidTokenURL(string(make([]byte, MaxTokenURLLength+1))),
		},
		{
			"case-description",
			NewMsgIssueToken("name", "coin", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", string(make([]byte, MaxTokenDescriptionLength+1)), TestIdentityString),
			ErrInvalidTokenDescription(string(make([]byte, MaxTokenDescriptionLength+1))),
		},
		{
			"case-identity1",
			NewMsgIssueToken("name", "coin", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", string(make([]byte, MaxTokenIdentityLength+1))),
			ErrInvalidTokenIdentity(string(make([]byte, MaxTokenIdentityLength+1))),
		},
		{
			"case-identity2",
			NewMsgIssueToken("name", "coin", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", ""),
			ErrNilTokenIdentity(),
		},
		{
			"nil name",
			NewMsgIssueToken("", "coin", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_ValidateBasic(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", testAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgTransferOwnership("123", testAddr, addr),
			ErrInvalidTokenSymbol("123"),
		},
		{
			"case-invalid1",
			NewMsgTransferOwnership("abc", sdk.AccAddress{}, testAddr),
			ErrNilTokenOwner(),
		},
		{
			"case-invalid2",
			NewMsgTransferOwnership("abc", testAddr, sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalid3",
			NewMsgTransferOwnership("abc", testAddr, testAddr),
			ErrTransferSelfTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgTransferOwnership.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgMintToken("abc", sdk.NewInt(10000), testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgMintToken("()2", sdk.NewInt(10000), testAddr),
			ErrInvalidTokenSymbol("()2"),
		},
		{
			"case-invalidOwner",
			NewMsgMintToken("abc", sdk.NewInt(10000), sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAmt",
			NewMsgMintToken("abc", sdk.NewInt(-1), testAddr),
			ErrInvalidTokenMintAmt(sdk.NewInt(-1).String()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgMintToken.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgBurnToken("abc", sdk.NewInt(10000), testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgBurnToken("w♞", sdk.NewInt(10000), testAddr),
			ErrInvalidTokenSymbol("w♞"),
		},
		{
			"case-invalidOwner",
			NewMsgBurnToken("abc", sdk.NewInt(10000), sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAmt",
			NewMsgBurnToken("abc", sdk.NewInt(-1), testAddr),
			ErrInvalidTokenBurnAmt(sdk.NewInt(-1).String()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBurnToken.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgForbidToken("abc", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidToken("*90", testAddr),
			ErrInvalidTokenSymbol("*90"),
		},
		{
			"case-invalidOwner",
			NewMsgForbidToken("abc", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgUnForbidToken("abc", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidToken("a¥0", testAddr),
			ErrInvalidTokenSymbol("a¥0"),
		},
		{
			"case-invalidOwner",
			NewMsgUnForbidToken("abc", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockAddrList()
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", testAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgAddTokenWhitelist("abcdefghi01234567", testAddr, whitelist),
			ErrInvalidTokenSymbol("abcdefghi01234567"),
		},
		{
			"case-invalidOwner",
			NewMsgAddTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidWhitelist",
			NewMsgAddTokenWhitelist("abc", testAddr, []sdk.AccAddress{}),
			ErrNilTokenWhitelist(),
		},
		{
			"case-nilWhitelist",
			NewMsgAddTokenWhitelist("abc", testAddr, []sdk.AccAddress{nilAddr, nilAddr}),
			ErrNilTokenWhitelist(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockAddrList()
	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", testAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgRemoveTokenWhitelist("a℃", testAddr, whitelist),
			ErrInvalidTokenSymbol("a℃"),
		},
		{
			"case-invalidOwner",
			NewMsgRemoveTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidWhitelist",
			NewMsgRemoveTokenWhitelist("abc", testAddr, []sdk.AccAddress{}),
			ErrNilTokenWhitelist(),
		},
		{
			"case-nilWhitelist",
			NewMsgRemoveTokenWhitelist("abc", testAddr, []sdk.AccAddress{nilAddr, nilAddr}),
			ErrNilTokenWhitelist(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_ValidateBasic(t *testing.T) {
	addresses := mockAddrList()
	tests := []struct {
		name string
		msg  MsgForbidAddr
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgForbidAddr("abc", testAddr, addresses),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidAddr("a⎝⎠", testAddr, addresses),
			ErrInvalidTokenSymbol("a⎝⎠"),
		},
		{
			"case-invalidOwner",
			NewMsgForbidAddr("abc", sdk.AccAddress{}, addresses),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAddr",
			NewMsgForbidAddr("abc", testAddr, []sdk.AccAddress{}),
			ErrNilForbiddenAddress(),
		},
		{
			"case-forbidSelf",
			NewMsgForbidAddr("abc", testAddr, []sdk.AccAddress{testAddr}),
			ErrTokenOwnerSelfForbidden(),
		},
		{
			"case-nilForbiddenAddress",
			NewMsgForbidAddr("abc", testAddr, []sdk.AccAddress{nilAddr, nilAddr}),
			ErrNilForbiddenAddress(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_ValidateBasic(t *testing.T) {
	addr := mockAddrList()
	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgUnForbidAddr("abc", testAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidAddr("a⥇", testAddr, addr),
			ErrInvalidTokenSymbol("a⥇"),
		},
		{
			"case-invalidOwner",
			NewMsgUnForbidAddr("abc", sdk.AccAddress{}, addr),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAddr",
			NewMsgUnForbidAddr("abc", testAddr, []sdk.AccAddress{}),
			ErrNilForbiddenAddress(),
		},
		{
			"case-nilForbiddenAddress",
			NewMsgUnForbidAddr("abc", testAddr, []sdk.AccAddress{nilAddr, nilAddr}),
			ErrNilForbiddenAddress(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgModifyTokenURL_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgModifyTokenInfo
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgModifyTokenInfo("abc", "www.abc.org", "abc example description", TestIdentityString, testAddr, "ABC token", "1000", "true", "true", "true", "true"),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgModifyTokenInfo("a😃", "www.abc.org", "abc example description", TestIdentityString, testAddr, "ABC token", "1000", "true", "true", "true", "true"),
			ErrInvalidTokenSymbol("a😃"),
		},
		{
			"case-invalidOwner",
			NewMsgModifyTokenInfo("abc", "www.abc.org", "abc example description", TestIdentityString, sdk.AccAddress{}, "ABC token", "1000", "true", "true", "true", "true"),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidURL",
			NewMsgModifyTokenInfo("abc", string(make([]byte, MaxTokenURLLength+1)), "abc example description", TestIdentityString, testAddr, "ABC token", "1000", "true", "true", "true", "true"),
			ErrInvalidTokenURL(string(make([]byte, MaxTokenURLLength+1))),
		},
		{
			"case-invalidDescription",
			NewMsgModifyTokenInfo("abc", "www.abc.org", string(make([]byte, MaxTokenDescriptionLength+1)), TestIdentityString, testAddr, "ABC token", "1000", "true", "true", "true", "true"),
			ErrInvalidTokenDescription(string(make([]byte, MaxTokenDescriptionLength+1))),
		},
		{
			"case-invalidIdentity",
			NewMsgModifyTokenInfo("abc", "www.abc.org", "abc example description", string(make([]byte, MaxTokenIdentityLength+1)), testAddr, "ABC token", "1000", "true", "true", "true", "true"),
			ErrInvalidTokenIdentity(string(make([]byte, MaxTokenIdentityLength+1))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgModifyTokenURL.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsg_Route(t *testing.T) {
	want := RouterKey
	tests := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			"issue-token",
			MsgIssueToken{},
		},
		{
			"transfer-ownership",
			MsgTransferOwnership{},
		},
		{
			"burn-token",
			MsgBurnToken{},
		},
		{
			"mint-token",
			MsgMintToken{},
		},
		{
			"forbid-token",
			MsgForbidToken{},
		},
		{
			"unforbid-token",
			MsgUnForbidToken{},
		},
		{
			"add_token_whitelist",
			MsgAddTokenWhitelist{},
		},
		{
			"remove-token-whitelist",
			MsgRemoveTokenWhitelist{},
		},
		{
			"forbid-addr",
			MsgForbidAddr{},
		},
		{
			"unforbid-addr",
			MsgUnForbidAddr{},
		},
		{
			"modify-token-info",
			MsgModifyTokenInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.Route(); got != want {
				t.Errorf("Msg.Route() = %v, want %v", got, want)
			}
		})
	}
}

func TestMsg_Type(t *testing.T) {
	tests := []struct {
		name string
		msg  sdk.Msg
		want string
	}{
		{
			"issue-token",
			MsgIssueToken{},
			"issue_token",
		},
		{
			"transfer-ownership",
			MsgTransferOwnership{},
			"transfer_ownership",
		},
		{
			"burn-token",
			MsgBurnToken{},
			"burn_token",
		},
		{
			"mint-token",
			MsgMintToken{},
			"mint_token",
		},
		{
			"forbid-token",
			MsgForbidToken{},
			"forbid_token",
		},
		{
			"unforbid-token",
			MsgUnForbidToken{},
			"unforbid_token",
		},
		{
			"add_token_whitelist",
			MsgAddTokenWhitelist{},
			"add_token_whitelist",
		},
		{
			"remove-token-whitelist",
			MsgRemoveTokenWhitelist{},
			"remove_token_whitelist",
		},
		{
			"forbid-addr",
			MsgForbidAddr{},
			"forbid_addr",
		},
		{
			"unforbid-addr",
			MsgUnForbidAddr{},
			"unforbid_addr",
		},
		{
			"modify-token-info",
			MsgModifyTokenInfo{},
			"modify_token_info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.Type(); got != tt.want {
				t.Errorf("Msg.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsg_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  sdk.Msg
		want []sdk.AccAddress
	}{
		{
			"issue-token",
			NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(10000), testAddr,
				false, false, false, false, "", "", TestIdentityString),
			[]sdk.AccAddress{testAddr},
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", testAddr, sdk.AccAddress{}),
			[]sdk.AccAddress{testAddr},
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", sdk.NewInt(10000), testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"mint-token",
			NewMsgMintToken("abc", sdk.NewInt(10000), testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"forbid-token",
			NewMsgForbidToken("abc", testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"unforbid-token",
			NewMsgUnForbidToken("abc", testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"remove-token-whitelist",
			NewMsgRemoveTokenWhitelist("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"forbid-addr",
			NewMsgForbidAddr("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"unforbid-addr",
			NewMsgUnForbidAddr("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"modify-token-url",
			NewMsgModifyTokenInfo("abc", "www.abc.com", "abc example description", TestIdentityString, testAddr,
				"ABC token", "1000", "true", "true", "true", "true"),
			[]sdk.AccAddress{testAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Msg.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsg_GetSignBytes(t *testing.T) {
	var owner, _ = sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")

	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addrList = []sdk.AccAddress{addr1, addr2}

	tests := []struct {
		name string
		msg  sdk.Msg
		want string
	}{
		{
			"issue-token",
			NewMsgIssueToken("ABC Token", "abc", sdk.NewInt(10000), owner,
				false, false, false, false, "", "", TestIdentityString),
			`{"type":"asset/MsgIssueToken","value":{"addr_forbiddable":false,"burnable":false,"description":"","identity":"552A83BA62F9B1F8","mintable":false,"name":"ABC Token","owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","token_forbiddable":false,"total_supply":"10000","url":""}}`,
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", owner, addr),
			`{"type":"asset/MsgTransferOwnership","value":{"new_owner":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","original_owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", sdk.NewInt(10000), owner),
			`{"type":"asset/MsgBurnToken","value":{"amount":"10000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"mint-token",
			NewMsgMintToken("abc", sdk.NewInt(10000), owner),
			`{"type":"asset/MsgMintToken","value":{"amount":"10000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"forbid-token",
			NewMsgForbidToken("abc", owner),
			`{"type":"asset/MsgForbidToken","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"unforbid-token",
			NewMsgUnForbidToken("abc", owner),
			`{"type":"asset/MsgUnForbidToken","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", owner, addrList),
			`{"type":"asset/MsgAddTokenWhitelist","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
		{
			"remove-token-whitelist",
			NewMsgRemoveTokenWhitelist("abc", owner, addrList),
			`{"type":"asset/MsgRemoveTokenWhitelist","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
		{
			"forbid-addr",
			NewMsgForbidAddr("abc", owner, addrList),
			`{"type":"asset/MsgForbidAddr","value":{"addresses":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"],"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"unforbid-addr",
			NewMsgUnForbidAddr("abc", owner, addrList),
			`{"type":"asset/MsgUnForbidAddr","value":{"addresses":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"],"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"modify-token-info",
			NewMsgModifyTokenInfo("abc", "www.abc.com", "abc example description", TestIdentityString, owner,
				"ABC token", "1000", "true", "true", "true", "true"),
			`{"type":"asset/MsgModifyTokenInfo","value":{"addr_forbiddable":"true","burnable":"true","description":"abc example description","identity":"552A83BA62F9B1F8","mintable":"true","name":"ABC token","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","token_forbiddable":"true","total_supply":"1000","url":"www.abc.com"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, string(tt.msg.GetSignBytes()))
		})
	}
}
