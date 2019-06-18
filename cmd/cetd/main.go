package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/coinexchain/dex/app"
	cet_init "github.com/coinexchain/dex/init"
	cet_server "github.com/coinexchain/dex/server"
	dex "github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
)

// cetd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	dex.InitSdkConfig()
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := createRootCmd(ctx, cdc)
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

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

func createRootCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:               "cetd",
		Short:             "CET Chain Daemon (server)",
		PersistentPreRunE: cet_server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(cet_init.InitCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.AddGenesisTokenCmd(ctx, cdc))
	rootCmd.AddCommand(cet_init.ValidateGenesisCmd(ctx, cdc))

	return rootCmd
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
