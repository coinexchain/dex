package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
)

var ResultParam *keeper.QueryAddrBalances
var ResultPath string

func TestQuery(t *testing.T) {
	cliutil.CliQuery = func(cdc *codec.Codec, path string, param interface{}) error {
		ResultParam = param.(*keeper.QueryAddrBalances)
		ResultPath = path
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetQueryCmd(nil)
	args := []string{
		"balances",
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/bankx/balances", ResultPath)
	assert.Equal(t, &keeper.QueryAddrBalances{Addr: addr}, ResultParam)

	args = []string{
		"balances",
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "decoding bech32 failed: checksum failed. Expected eqv7uv, got agaytu.", err.Error())
}

func TestQueryParams(t *testing.T) {
	cliutil.CliQuery = func(cdc *codec.Codec, path string, param interface{}) error {
		ResultParam = nil
		ResultPath = path
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetQueryCmd(nil)
	args := []string{
		"params",
	}
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/bankx/parameters", ResultPath)
}
