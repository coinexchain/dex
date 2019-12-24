package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	dex.InitSdkConfig()

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	rootCmd := createRootCmd(cdc)

	executor := cli.PrepareMainCmd(rootCmd, "", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func createRootCmd(cdc *codec.Codec) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cetdev",
		Short: "Command line tools for CoinEx Chain developers",
	}

	rootCmd.AddCommand(
		ExampleGenesisCmd(cdc),
		DefaultParamsCmd(),
		CosmosHubParamsCmd(cdc),
		RestEndpointsCmd(registerRoutes),
		//ShowCommandTreeCmd(),
	)

	return rootCmd
}

func registerRoutes(rs *lcd.RestServer) {
	//registerSwaggerUI(rs)
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}
