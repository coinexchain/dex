package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	types "github.com/coinexchain/dex/modules/alias/internal/types"
)

var ResultMsg *types.MsgAliasUpdate

func CliRunCommandForTest(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
	ResultMsg = msg.(*types.MsgAliasUpdate)
	return nil
}

func TestCmd(t *testing.T) {
	cliutil.CliRunCommand = CliRunCommandForTest

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetTxCmd(nil)

	args := []string{
		"add",
		"super_super_boy",
		"--as-default=yes",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgAliasUpdate{
		Alias:     "super_super_boy",
		IsAdd:     true,
		AsDefault: true,
	}
	assert.Equal(t, msg, ResultMsg)

	args = []string{
		"add",
		"super_boy",
		"--as-default=no",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	msg = &types.MsgAliasUpdate{
		Alias:     "super_boy",
		IsAdd:     true,
		AsDefault: false,
	}
	assert.Equal(t, msg, ResultMsg)

	args = []string{
		"add",
		"super_boy",
		"--as-default=ok",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "Invalid value for --as-default, only 'yes' and 'no' are valid", err.Error())

	args = []string{
		"remove",
		"super_boy",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	msg = &types.MsgAliasUpdate{
		Alias: "super_boy",
		IsAdd: false,
	}
	assert.Equal(t, msg, ResultMsg)
}
