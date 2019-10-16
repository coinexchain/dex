package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	types "github.com/coinexchain/dex/modules/alias/internal/types"
)

var ResultMsg *types.MsgAliasUpdate

func CliRunCommandForTest(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg.SetAccAddress(senderAddr)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	ResultMsg = msg.(*types.MsgAliasUpdate)
	return nil
}

func TestCmd(t *testing.T) {
	cliutil.CliRunCommand = CliRunCommandForTest

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetTxCmd(nil)

	addr, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234abcde")
	addrStr := addr.String()

	args := []string{
		"add",
		"super_super_boy",
		"--as-default=yes",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgAliasUpdate{
		Owner:     addr,
		Alias:     "super_super_boy",
		IsAdd:     true,
		AsDefault: true,
	}
	assert.Equal(t, msg, ResultMsg)

	args = []string{
		"add",
		"super_boy",
		"--as-default=no",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	msg = &types.MsgAliasUpdate{
		Owner:     addr,
		Alias:     "super_boy",
		IsAdd:     true,
		AsDefault: false,
	}
	assert.Equal(t, msg, ResultMsg)

	args = []string{
		"add",
		"super_boy",
		"--as-default=ok",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Error(t, err)
	assert.Equal(t, "Invalid value for --as-default, only 'yes' and 'no' are valid", err.Error())

	args = []string{
		"remove",
		"super_boy",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	msg = &types.MsgAliasUpdate{
		Owner: addr,
		Alias: "super_boy",
		IsAdd: false,
	}
	assert.Equal(t, msg, ResultMsg)
}
