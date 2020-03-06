package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	staking_cli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	authxcmd "github.com/coinexchain/cet-sdk/modules/authx/client/cli"
	bankxcmd "github.com/coinexchain/cet-sdk/modules/bankx/client/cli"
	distrxcmd "github.com/coinexchain/cet-sdk/modules/distributionx/client/cli"
	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
	_ "github.com/coinexchain/dex/cmd/cetcli/statik"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	dex.InitSdkConfig()

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()
	rootCmd := createRootCmd(cdc)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func createRootCmd(cdc *codec.Codec) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cetcli",
		Short: "Command line interface for interacting with cetd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(FlagSwaggerHost, "", "Default host of swagger API")
	rootCmd.PersistentFlags().Bool(FlagDefaultHTTP, false, "Use Http as default schema")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.Cmd,
		client.NewCompletionCmd(rootCmd, true),
		client.LineBreak,
	)

	config := sdk.GetConfig()
	// fix `keys parse` cmd
	keys.Bech32Prefixes = []string{
		config.GetBech32AccountAddrPrefix(),
		config.GetBech32AccountPubPrefix(),
		config.GetBech32ValidatorAddrPrefix(),
		config.GetBech32ValidatorPubPrefix(),
		config.GetBech32ConsensusAddrPrefix(),
		config.GetBech32ConsensusPubPrefix(),
	}

	fixDescriptions(rootCmd)
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
	if err := viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag)); err != nil {
		return err
	}

	return bindSwaggerFlags(cmd)
}

func queryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authxcmd.GetAccountXCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankxcmd.SendTxCmd(cdc),
		bankxcmd.RequireMemoCmd(cdc),
		distrxcmd.DonateTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	fixUnknownFlagIssue(txCmd)

	return txCmd
}

func fixUnknownFlagIssue(txCmd *cobra.Command) {
	cmd, _, err := txCmd.Find([]string{"staking", "edit-validator"})
	if err == nil {
		cmd.Flags().AddFlagSet(staking_cli.FsMinSelfDelegation)
	}
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	registerSwaggerUI(rs)
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func fixDescriptions(cmd *cobra.Command) {
	// cosmosvalconspubXXX -> coinexvalconspubXXX
	if idx := strings.Index(cmd.Long, "cosmosvalconspub"); idx >= 0 {
		cosmosPubKey := cmd.Long[idx : idx+83]
		pubKey := consPubKeyFromBech32("cosmosvalconspub", cosmosPubKey)
		dexPubKey := sdk.MustBech32ifyConsPub(pubKey)
		cmd.Long = strings.Replace(cmd.Long, cosmosPubKey, dexPubKey, -1)
		//fmt.Printf("%s -> %s\n", cosmosPubKey, dexPubKey)
	}
	// cosmosvaloperXXX -> coinexvaloperXXX
	for {
		if idx := strings.Index(cmd.Long, "cosmosvaloper"); idx >= 0 {
			cosmosValAddr := cmd.Long[idx : idx+52]
			rawValAddr := rawAddressFromBech32("cosmosvaloper", cosmosValAddr)
			dexValAddr := sdk.ValAddress(rawValAddr).String()
			cmd.Long = strings.Replace(cmd.Long, cosmosValAddr, dexValAddr, -1)
			//fmt.Printf("%s -> %s\n", cosmosValAddr, dexValAddr)
		} else {
			break
		}
	}
	// cosmosXXX -> coinexXXX
	if idx := strings.Index(cmd.Long, "cosmos1..."); idx >= 0 {
		cmd.Long = strings.ReplaceAll(cmd.Long, "cosmos1...", "coinex1...")
	}
	if idx := strings.Index(cmd.Long, "cosmos1"); idx >= 0 {
		cosmosAccAddr := cmd.Long[idx : idx+45]
		rawAccAddr := rawAddressFromBech32("cosmos", cosmosAccAddr)
		dexAccAddr := sdk.AccAddress(rawAccAddr).String()
		cmd.Long = strings.Replace(cmd.Long, cosmosAccAddr, dexAccAddr, -1)
		//fmt.Printf("%s -> %s\n", cosmosAccAddr, dexAccAddr)
	}
	if idx := strings.Index(cmd.Long, "0stake"); idx >= 0 {
		cmd.Long = strings.Replace(cmd.Long, "0stake", "0cet", -1)
		//fmt.Printf("%s -> %s\n", "stake", dex.CET)
	}

	// uatom -> cet
	for _, flagName := range []string{client.FlagFees, client.FlagGasPrices} {
		if flag := cmd.Flag(flagName); flag != nil {
			flag.Usage = strings.Replace(flag.Usage, "uatom", dex.CET, -1)
			//fmt.Printf("%s -> %s\n", "uatom", dex.CET)
		}
	}

	if cmd.Name() == "submit-proposal" {
		cmd.Long = strings.ReplaceAll(cmd.Long, `10test`, `10cet`)
	}
	if cmd.Name() == "param-change" {
		cmd.Long = strings.ReplaceAll(cmd.Long, `"stake"`, `"cet"`)
	}
	if cmd.Name() == "community-pool-spend" {
		cmd.Long = strings.ReplaceAll(cmd.Long, `"stake"`, `"cet"`)
		cmd.Long = strings.ReplaceAll(cmd.Long, "Pay me some Atoms!", "Pay me some CETs!")
	}

	if len(cmd.Commands()) > 0 {
		for _, subCmd := range cmd.Commands() {
			fixDescriptions(subCmd)
		}
	}
}

func rawAddressFromBech32(prefix, address string) (addr []byte) {
	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		panic(err)
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		panic(err)
	}

	return bz
}

func consPubKeyFromBech32(prefix, pubkey string) (pk crypto.PubKey) {
	bz, err := sdk.GetFromBech32(pubkey, prefix)
	if err != nil {
		panic(err)
	}

	pk, err = cryptoAmino.PubKeyFromBytes(bz)
	if err != nil {
		panic(err)
	}

	return pk
}
