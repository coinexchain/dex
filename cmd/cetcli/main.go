package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	crisisclient "github.com/cosmos/cosmos-sdk/x/crisis/client"
	distcmd "github.com/cosmos/cosmos-sdk/x/distribution"
	distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
	gv "github.com/cosmos/cosmos-sdk/x/gov"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	gov "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	sl "github.com/cosmos/cosmos-sdk/x/slashing"
	slashingclient "github.com/cosmos/cosmos-sdk/x/slashing/client"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	st "github.com/cosmos/cosmos-sdk/x/staking"
	stakingclient "github.com/cosmos/cosmos-sdk/x/staking/client"
	staking "github.com/cosmos/cosmos-sdk/x/staking/client/rest"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/cmd/cetcli/dev"
	_ "github.com/coinexchain/dex/cmd/cetcli/statik"
	as "github.com/coinexchain/dex/modules/asset"
	assclient "github.com/coinexchain/dex/modules/asset/client"
	assrest "github.com/coinexchain/dex/modules/asset/client/rest"
	"github.com/coinexchain/dex/modules/authx"
	authxcmd "github.com/coinexchain/dex/modules/authx/client/cli"
	authxrest "github.com/coinexchain/dex/modules/authx/client/rest"
	bankxcmd "github.com/coinexchain/dex/modules/bankx/client/cli"
	bankxrest "github.com/coinexchain/dex/modules/bankx/client/rest"
	mktclient "github.com/coinexchain/dex/modules/market/client"
	mktrest "github.com/coinexchain/dex/modules/market/client/rest"
	dex "github.com/coinexchain/dex/types"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	dex.InitSdkConfig()

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	rootCmd := createRootCmd(cdc)
	fixDescriptions(rootCmd)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func createRootCmd(cdc *amino.Codec) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cetcli",
		Short: "Command line interface for interacting with cetd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Module clients hold cli commands (tx,query) and lcd routes
	// TODO: Make the lcd command take a list of ModuleClient
	mc := []sdk.ModuleClients{
		assclient.NewModuleClient(as.StoreKey, cdc),
		mktclient.NewModuleClient(as.StoreKey, cdc),
		govClient.NewModuleClient(gv.StoreKey, cdc),
		distClient.NewModuleClient(distcmd.StoreKey, cdc),
		stakingclient.NewModuleClient(st.StoreKey, cdc),
		slashingclient.NewModuleClient(sl.StoreKey, cdc),
		crisisclient.NewModuleClient(sl.StoreKey, cdc),
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		configCmd(),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
		client.NewCompletionCmd(rootCmd, true),
		client.LineBreak,
		dev.DevCmd(cdc),
	)

	return rootCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

func configCmd() *cobra.Command {
	cmd := client.ConfigCmd(app.DefaultCLIHome)
	cmd.Short = "Create or query a CoinEx Chain CLI configuration file"
	return cmd
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authxcmd.GetAccountXCmd(authx.StoreKey, cdc),
	)

	for _, m := range mc {
		mQueryCmd := m.GetQueryCmd()
		if mQueryCmd != nil {
			queryCmd.AddCommand(mQueryCmd)
		}
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankxcmd.SendTxCmd(cdc),
		bankxcmd.RequireMemoCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		tx.GetBroadcastCommand(cdc),
		tx.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	registerSwaggerUI(rs)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	authxrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, authx.StoreKey)
	bankxrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, distcmd.StoreKey)
	staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	slashing.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	gov.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	assrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, as.StoreKey)
	mktrest.RegisterTXRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
}

func registerSwaggerUI(rs *lcd.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", staticServer))
}

func fixDescriptions(cmd *cobra.Command) {
	gaiacli := "gaiacli"
	cetcli := "cetcli"

	cmd.Short = strings.Replace(cmd.Short, gaiacli, cetcli, -1)
	cmd.Long = strings.Replace(cmd.Long, gaiacli, cetcli, -1)
	if flagFee := cmd.Flag(client.FlagFees); flagFee != nil {
		flagFee.Usage = "Fees to pay along with transaction; eg: 100cet"
	}
	if flagGas := cmd.Flag(client.FlagGasPrices); flagGas != nil {
		flagGas.Usage = "Gas prices to determine the transaction fee (e.g. 100cet)"
	}

	if len(cmd.Commands()) > 0 {
		for _, subCmd := range cmd.Commands() {
			fixDescriptions(subCmd)
		}
	}
}
