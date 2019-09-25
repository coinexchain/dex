package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

func RequireMemoCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "require-memo <bool>",
		Short: "Mark if memo is required to receive coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			required, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgSetTransferMemoRequired(nil, required)
			return cliutil.CliRunCommand(cdc, &msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
