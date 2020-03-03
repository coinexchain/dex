package main

import (
	"encoding/json"
	"io"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	assetcli "github.com/coinexchain/cet-sdk/modules/asset/client/cli"
	"github.com/coinexchain/cet-sdk/msgqueue"
	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/app/plugin"
)

// cetd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	plugin.SetReloadPluginSignal(syscall.SIGUSR1)
	msgqueue.SetMkFifoFunc(syscall.Mkfifo)

	dex.InitSdkConfig()
	rootCmd := createCetdCmd()

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func createCetdCmd() *cobra.Command {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "cetd",
		Short:             "CoinEx Chain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	addInitCommands(ctx, cdc, rootCmd)
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")

	return rootCmd
}

func addInitCommands(ctx *server.Context, cdc *codec.Codec, rootCmd *cobra.Command) {
	rawBasicManager := app.ModuleBasics.BasicManager

	initCmd := genutilcli.InitCmd(ctx, cdc, rawBasicManager, app.DefaultNodeHome)
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		adjustBlockCommitSpeed(ctx.Config)
	}
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, rawBasicManager, staking.AppModuleBasic{},
		genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, rawBasicManager))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(assetcli.AddGenesisTokenCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(testnetCmd(ctx, cdc, app.ModuleBasics, genaccounts.AppModuleBasic{}))
	rootCmd.AddCommand(migrateCmd(cdc))
}

func adjustBlockCommitSpeed(config *tmconfig.Config) {
	c := config.Consensus
	c.TimeoutCommit = 2000 * time.Millisecond
	c.PeerGossipSleepDuration = 20 * time.Millisecond
	c.PeerQueryMaj23SleepDuration = 100 * time.Millisecond
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	cetChainApp := app.NewCetChainApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetCheckTxWithMsgHandle(viper.GetBool(server.FlagCheckTxWithMsgHandle)),
	)
	checkMinGasPrice(cetChainApp, logger)
	return cetChainApp
}

func checkMinGasPrice(bApp *app.CetChainApp, logger log.Logger) {
	ctx := bApp.NewContext(true, abci.Header{})
	minGasPrice := ctx.MinGasPrices().AmountOf(dex.CET)
	if !minGasPrice.IsPositive() {
		panic("--minimum-gas-prices option not set!")
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
