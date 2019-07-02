package dev

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
)

// add server commands
func DevCmd(cdc *codec.Codec, registerRoutesFn func(*lcd.RestServer)) *cobra.Command {
	devCmd := &cobra.Command{
		Use:   "dev",
		Short: "Dev subcommands",
	}

	devCmd.AddCommand(
		ExampleGenesisCmd(cdc),
		TestnetGenesisCmd(cdc),
		DefaultParamsCmd(cdc),
		CosmosHubParamsCmd(),
		ShowCommandTreeCmd(),
		RestEndpointsCmd(cdc, registerRoutesFn),
	)

	return devCmd
}
