package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
)

var ResultParam *keepers.QueryAliasInfoParam
var ResultPath string

func CliQueryForTest(_ *codec.Codec, path string, param interface{}) error {
	ResultParam = param.(*keepers.QueryAliasInfoParam)
	ResultPath = path
	return nil
}

func TestQuery(t *testing.T) {
	cliutil.CliQuery = CliQueryForTest

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetQueryCmd(nil)

	args := []string{
		"address-of-alias",
		"super_super_boy",
	}
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "custom/alias/alias-info", ResultPath)
	assert.Equal(t, &keepers.QueryAliasInfoParam{Alias: "super_super_boy", QueryOp: keepers.GetAddressFromAlias}, ResultParam)

	args = []string{
		"aliases-of-address",
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "custom/alias/alias-info", ResultPath)
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	assert.Equal(t, &keepers.QueryAliasInfoParam{
		Owner:   addr,
		QueryOp: keepers.ListAliasOfAccount,
	}, ResultParam)

	args = []string{
		"aliases-of-address",
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "decoding bech32 failed: checksum failed. Expected eqv7uv, got agaytu.", err.Error())
}
