package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func TestSendTxCmd(t *testing.T) {
	var resultMsg *types.MsgSend
	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		resultMsg = msg.(*types.MsgSend)
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := SendTxCmd(nil)

	args := []string{
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
		"1000000000cet",
		"--from=bob",
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	amount := dex.NewCetCoins(1000000000)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgSend{
		ToAddress:  addr,
		Amount:     amount,
		UnlockTime: 0,
	}
	assert.Equal(t, msg, resultMsg)
}

func TestRequireMemoCmd(t *testing.T) {
	var resultMsg *types.MsgSetMemoRequired
	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		resultMsg = msg.(*types.MsgSetMemoRequired)
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := RequireMemoCmd(nil)

	args := []string{
		"true",
		"--from=bob",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgSetMemoRequired{
		Required: true,
	}
	assert.Equal(t, msg, resultMsg)
}
