package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	staking_cli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	stakingQueryCmd := staking_cli.GetQueryCmd(staking.QuerierRoute, cdc)
	BondPoolCmd := client.GetCommands(GetCmdQueryPool(cdc))[0]

	//replace pool cmd with new bondPoolCmd which can also show the non-bondable-cet-tokens in locked positions
	return replacePoolCmd(stakingQueryCmd, BondPoolCmd)
}

func replacePoolCmd(stakingQueryCmd *cobra.Command, BondPoolCmd *cobra.Command) *cobra.Command {
	var oldPoolCmd *cobra.Command
	for _, cmd := range stakingQueryCmd.Commands() {
		if cmd.Use == "pool" {
			oldPoolCmd = cmd
		}
	}

	stakingQueryCmd.RemoveCommand(oldPoolCmd)
	stakingQueryCmd.AddCommand(BondPoolCmd)

	return stakingQueryCmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current staking pool values",
		Long: strings.TrimSpace(`Query values for amounts stored in the staking pool:

$ cetcli query staking pool
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData("custom/stakingx/pool", nil)
			if err != nil {
				return err
			}

			println(string(res))
			return nil
		},
	}
}
