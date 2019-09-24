package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/distributionx/types"
)

// DonateTxCmd will create a DonateToCommunityPool tx and sign it with the given key.
func DonateTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate [amount]",
		Short: "Donate to community pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDonateToCommunityPool(nil, coins)
			return cliutil.CliRunCommand(cdc, &msg)

		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.MarkFlagRequired(client.FlagFrom)
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	return cmd
}
