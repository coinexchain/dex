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

	gaia_init "github.com/cosmos/cosmos-sdk/cmd/gaia/init"

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
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return genGenesisJSON(ctx, cdc, args[0])
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

func genGenesisJSON(ctx *server.Context, cdc *codec.Codec, moniker string) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))
	config.Moniker = moniker

	chainID := viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
	}

	// generate node_key.json & priv_validator_key.json
	nodeID, _, err := gaia_init.InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	var appState json.RawMessage
	genFile := config.GenesisFile()

	// generate genesis.json
	if appState, err = initializeEmptyGenesis(cdc, genFile, viper.GetBool(flagOverwrite)); err != nil {
		return err
	}
	if err = gaia_init.ExportGenesisFile(genFile, chainID, nil, appState); err != nil {
		return err
	}

	// generate config.toml
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
	return displayInfo(cdc, toPrint)
}

func initializeEmptyGenesis(
	cdc *codec.Codec, genFile string, overwrite bool,
) (appState json.RawMessage, err error) {

	if !overwrite && common.FileExists(genFile) {
		return nil, fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	return codec.MarshalJSONIndent(cdc, app.NewDefaultGenesisState())
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
