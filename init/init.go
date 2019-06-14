package init

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
)

const (
	flagOverwrite    = "overwrite"
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message"`
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return initFn(ctx, cdc, args[0])
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

func initFn(ctx *server.Context, cdc *codec.Codec, moniker string) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))
	config.Moniker = moniker

	// generate node_key.json & priv_validator_key.json
	nodeID, _, err := InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	// generate genesis.json
	chainID, appState, err := initializeGenesisFile(cdc, config.GenesisFile())
	if err != nil {
		return err
	}

	// generate config.toml
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
	return displayInfo(cdc, toPrint)
}

func initializeGenesisFile(cdc *codec.Codec, genFile string) (chainID string, appState json.RawMessage, err error) {
	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
	}

	if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {
		err = fmt.Errorf("genesis.json file already exists: %v", genFile)
		return
	}

	if appState, err = codec.MarshalJSONIndent(cdc, app.NewDefaultGenesisState()); err != nil {
		return
	}

	err = ExportGenesisFile(genFile, chainID, nil, appState)
	return
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", string(out)) // nolint: errcheck
	return nil
}
