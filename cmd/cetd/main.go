package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/app"
	assetcli "github.com/coinexchain/dex/modules/asset/client/cli"
	dexserver "github.com/coinexchain/dex/server"
	dex "github.com/coinexchain/dex/types"
)

// cetd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	dex.InitSdkConfig()
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := createCetdCmd(ctx, cdc)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func createCetdCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:               "cetd",
		Short:             "CoinEx Chain Daemon (server)",
		PersistentPreRunE: dexserver.PersistentPreRunEFn(ctx),
	}

	addInitCommands(ctx, cdc, rootCmd)
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	rootCmd.AddCommand(version.Cmd)

	return rootCmd
}

func addInitCommands(ctx *server.Context, cdc *amino.Codec, rootCmd *cobra.Command) {
	initCmd := genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome)
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		adjustBlockCommitSpeed(ctx.Config)
	}
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
		genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(assetcli.AddGenesisTokenCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(testnetCmd(ctx, cdc))
}

func adjustBlockCommitSpeed(config *tmconfig.Config) {
	c := config.Consensus
	c.TimeoutCommit = 2100 * time.Millisecond
	c.PeerGossipSleepDuration = 20 * time.Millisecond
	c.PeerQueryMaj23SleepDuration = 100 * time.Millisecond
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	cetChainApp := app.NewCetChainApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
	checkMinGasPrice(cetChainApp, logger)
	return cetChainApp
}

func checkMinGasPrice(bApp *app.CetChainApp, logger log.Logger) {
	ctx := bApp.NewContext(true, abci.Header{})
	minGasPrice := ctx.MinGasPrices().AmountOf(dex.CET)
	if !minGasPrice.IsPositive() {
		logger.Info("--minimum-gas-prices option not set!")
	}
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		gApp := app.NewCetChainApp(logger, db, traceStore, false, uint(1))
		err := gApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	gApp := app.NewCetChainApp(logger, db, traceStore, true, uint(1))
	return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
