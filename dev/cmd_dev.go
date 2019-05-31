package dev

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
)

// add server commands
func DevCmd(cdc *codec.Codec) *cobra.Command {
	devCmd := &cobra.Command{
		Use:   "dev",
		Short: "Dev subcommands",
	}

	devCmd.AddCommand(
		ExampleGenesisCmd(cdc),
	)

	return devCmd
}
