package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func TestSendTxCmd(t *testing.T) {
	var resultMsg *types.MsgSend
	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		cliCtx := context.NewCLIContext().WithCodec(cdc)
		senderAddr := cliCtx.GetFromAddress()
		msg.SetAccAddress(senderAddr)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		resultMsg = msg.(*types.MsgSend)
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := SendTxCmd(nil)

	addr, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234abcde")
	addrStr := addr.String()

	args := []string{
		"coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
		"1000000000cet",
		"--from=" + addrStr,
		"--generate-only",
	}
	addr1, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	amount := dex.NewCetCoins(1000000000)
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgSend{
		FromAddress: addr,
		ToAddress:   addr1,
		Amount:      amount,
		UnlockTime:  0,
	}
	assert.Equal(t, msg, resultMsg)
}

func TestRequireMemoCmd(t *testing.T) {
	var resultMsg *types.MsgSetMemoRequired
	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		cliCtx := context.NewCLIContext().WithCodec(cdc)
		senderAddr := cliCtx.GetFromAddress()
		msg.SetAccAddress(senderAddr)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		resultMsg = msg.(*types.MsgSetMemoRequired)
		return nil
	}

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := RequireMemoCmd(nil)

	addr, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234abcde")
	addrStr := addr.String()

	args := []string{
		"true",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	msg := &types.MsgSetMemoRequired{
		Address:  addr,
		Required: true,
	}
	assert.Equal(t, msg, resultMsg)
}
