package asset

import (
	"os"
	"reflect"
	"testing"

	asset_types "github.com/coinexchain/dex/modules/asset/types"
	dex "github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

func TestMsgIssueToken_ValidateBasic(t *testing.T) {
	invalidTokenSymbol := ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}")

	tests := []struct {
		name string
		msg  MsgIssueToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, tAccAddr,
				false, false, false, false, "", ""),
			nil,
		},
		{
			"case-name",
			NewMsgIssueToken(string(make([]byte, 32+1)), "abc", 100000, tAccAddr,
				false, false, false, false, "", ""),
			ErrorInvalidTokenName("token name is limited to 32 unicode characters"),
		},
		{
			"case-owner",
			NewMsgIssueToken("ABC Token", "abc", 100000, sdk.AccAddress{},
				false, false, false, false, "", ""),
			ErrorInvalidTokenOwner("token owner is invalid"),
		},
		{
			"case-symbol1",
			NewMsgIssueToken("ABC Token", "1aa", 100000, tAccAddr,
				false, false, false, false, "", ""),
			invalidTokenSymbol,
		},
		{
			"case-symbol2",
			NewMsgIssueToken("ABC Token", "A999", 100000, tAccAddr,
				false, false, false, false, "", ""),
			invalidTokenSymbol,
		},
		{
			"case-symbol3",
			NewMsgIssueToken("ABC Token", "aa1234567", 100000, tAccAddr,
				false, false, false, false, "", ""),
			invalidTokenSymbol,
		},
		{
			"case-symbol4",
			NewMsgIssueToken("ABC Token", "a*aa", 100000, tAccAddr,
				false, false, false, false, "", ""),
			invalidTokenSymbol,
		},
		{
			"case-totalSupply1",
			NewMsgIssueToken("ABC Token", "abc", 9E18+1, tAccAddr,
				false, false, false, false, "", ""),
			ErrorInvalidTokenSupply("token total supply before 1e8 boosting should be less than 90 billion"),
		},
		{
			"case-totalSupply2",
			NewMsgIssueToken("ABC Token", "abc", -1, tAccAddr,
				false, false, false, false, "", ""),
			ErrorInvalidTokenSupply("token total supply must be positive"),
		},
		{
			"case-url",
			NewMsgIssueToken("name", "coin", 2100, tAccAddr,
				false, false, false, false, string(make([]byte, 100+1)), ""),
			ErrorInvalidTokenURL("token url is limited to 100 unicode characters"),
		},
		{
			"case-description",
			NewMsgIssueToken("name", "coin", 2100, tAccAddr,
				false, false, false, false, "", string(make([]byte, 1024+1))),
			ErrorInvalidTokenDescription("token description is limited to 1k size"),
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
			NewMsgTransferOwnership("abc", tAccAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgTransferOwnership("123", tAccAddr, addr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
		{
			"case-invalid3",
			NewMsgTransferOwnership("abc", tAccAddr, tAccAddr),
			ErrorInvalidTokenOwner("Can not and no need to transfer ownership to self"),
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
			NewMsgMintToken("abc", 10000, tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgMintToken("()2", 10000, tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			NewMsgBurnToken("abc", 10000, tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgBurnToken("w♞", 10000, tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			NewMsgForbidToken("abc", tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidToken("*90", tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-invalidOwner",
			NewMsgForbidToken("abc", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("forbid token need a valid owner addr"),
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
			NewMsgUnForbidToken("abc", tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidToken("a¥0", tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-invalidOwner",
			NewMsgUnForbidToken("abc", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("forbid token need a valid owner addr"),
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
	whitelist := mockWhitelist()
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", tAccAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgAddTokenWhitelist("abcdefghi", tAccAddr, whitelist),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", tAccAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgRemoveTokenWhitelist("a℃", tAccAddr, whitelist),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgForbidAddr("abc", tAccAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidAddr("a⎝⎠", tAccAddr, addr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
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
			"base-case",
			NewMsgUnForbidAddr("abc", tAccAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidAddr("a⥇", tAccAddr, addr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
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
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgModifyTokenURL_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgModifyTokenURL
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgModifyTokenURL("abc", "www.abc.org", tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgModifyTokenURL("a😃", "www.abc.org", tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-invalidOwner",
			NewMsgModifyTokenURL("abc", "www.abc.org", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("modify token url need a valid owner addr"),
		},
		{
			"case-invalidURL",
			NewMsgModifyTokenURL("abc", string(make([]byte, 100+1)), tAccAddr),
			ErrorInvalidTokenURL("token url is limited to 100 unicode characters"),
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
func TestMsgModifyTokenDescription_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgModifyTokenDescription
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgModifyTokenDescription("abc", "abc example description", tAccAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgModifyTokenDescription("a❡", "abc example description", tAccAddr),
			ErrorInvalidTokenSymbol("token symbol not match with [a-z][a-z0-9]{1,7}"),
		},
		{
			"case-invalidOwner",
			NewMsgModifyTokenDescription("abc", "abc example description", sdk.AccAddress{}),
			ErrorInvalidTokenOwner("modify token description need a valid owner addr"),
		},
		{
			"case-invalidDescription",
			NewMsgModifyTokenDescription("abc", string(make([]byte, 1024+1)), tAccAddr),
			ErrorInvalidTokenDescription("token description is limited to 1k size"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMsgModifyTokenDescription.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsg_Route(t *testing.T) {
	want := asset_types.RouterKey
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
			"modify-token-url",
			MsgModifyTokenURL{},
		},
		{
			"modify-token-description",
			MsgModifyTokenDescription{},
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
			"modify-token-url",
			MsgModifyTokenURL{},
			"modify_token_url",
		},
		{
			"modify-token-description",
			MsgModifyTokenDescription{},
			"modify_token_description",
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
			NewMsgIssueToken("ABC Token", "abc", 100000, tAccAddr,
				false, false, false, false, "", ""),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", tAccAddr, sdk.AccAddress{}),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", 100000, tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"mint-token",
			NewMsgMintToken("abc", 100000, tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"forbid-token",
			NewMsgForbidToken("abc", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"unforbid-token",
			NewMsgUnForbidToken("abc", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", tAccAddr, mockWhitelist()),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"remove-token-whitelist",
			NewMsgRemoveTokenWhitelist("abc", tAccAddr, mockWhitelist()),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"forbid-addr",
			NewMsgForbidAddr("abc", tAccAddr, mockAddresses()),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"unforbid-addr",
			NewMsgUnForbidAddr("abc", tAccAddr, mockAddresses()),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"modify-token-url",
			NewMsgModifyTokenURL("abc", "www.abc.com", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
		},
		{
			"modify-token-description",
			NewMsgModifyTokenDescription("abc", "abc example description", tAccAddr),
			[]sdk.AccAddress{tAccAddr},
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
			NewMsgIssueToken("ABC Token", "abc", 100000, owner,
				false, false, false, false, "", ""),
			`{"type":"asset/MsgIssueToken","value":{"addr_forbiddable":false,"burnable":false,"description":"","mintable":false,"name":"ABC Token","owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","token_forbiddable":false,"total_supply":"100000","url":""}}`,
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", owner, addr),
			`{"type":"asset/MsgTransferOwnership","value":{"new_owner":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","original_owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", 100000, owner),
			`{"type":"asset/MsgBurnToken","value":{"amount":"100000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"mint-token",
			NewMsgMintToken("abc", 100000, owner),
			`{"type":"asset/MsgMintToken","value":{"amount":"100000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
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
			"modify-token-url",
			NewMsgModifyTokenURL("abc", "www.abc.com", owner),
			`{"type":"asset/MsgModifyTokenURL","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","url":"www.abc.com"}}`,
		},
		{
			"modify-token-description",
			NewMsgModifyTokenDescription("abc", "abc example description", owner),
			`{"type":"asset/MsgModifyTokenDescription","value":{"description":"abc example description","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Msg.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}
