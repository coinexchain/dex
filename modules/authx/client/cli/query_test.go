package cli

import (
	"github.com/coinexchain/dex/client/cliutil"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ResultParam *auth.QueryAccountParams
var ResultPath string

func TestQuery(t *testing.T) {
	cliutil.CliQuery = func(cdc *codec.Codec, path string, param interface{}) error {
		ResultParam = param.(*auth.QueryAccountParams)
		ResultPath = path
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetAccountXCmd(nil)
	args := []string{
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/accx/accountMix", ResultPath)
	assert.Equal(t, &auth.QueryAccountParams{Address: addr}, ResultParam)

	args = []string{
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
	assert.Equal(t, "custom/accx/parameters", ResultPath)
}
