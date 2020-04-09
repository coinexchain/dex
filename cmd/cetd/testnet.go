package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/stakingx"
	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
)

const nodeDirPerm = 0755

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagNodeCLIHome       = "node-cli-home"
	flagStartingIPAddress = "starting-ip-address"

	testnetTokenSupply       = sdk.NewInt(588788547005740000)
	testnetMinSelfDelegation = int64(10000e8)

	integrationTestChainID = "coinex-integrationtest"
)

type testnetNodeInfo struct {
	nodeID    string
	valPubKey crypto.PubKey
	acc       genaccounts.GenesisAccount
	genFile   string
}

type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis map[string]json.RawMessage,
		iterateFn func(authexported.Account) (stop bool),
	)
}

// get cmd to initialize all files for tendermint testnet and application
func testnetCmd(ctx *server.Context, cdc *codec.Codec,
	mbm dex.OrderedBasicManager, genAccIterator GenesisAccountsIterator) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a Cetd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	cetd testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config := ctx.Config

			outputDir := viper.GetString(flagOutputDir)
			chainID := viper.GetString(client.FlagChainID)
			minGasPrices := viper.GetString(server.FlagMinGasPrices)
			nodeDirPrefix := viper.GetString(flagNodeDirPrefix)
			nodeDaemonHome := viper.GetString(flagNodeDaemonHome)
			nodeCLIHome := viper.GetString(flagNodeCLIHome)
			startingIPAddress := viper.GetString(flagStartingIPAddress)
			numValidators := viper.GetInt(flagNumValidators)

			return initTestnet(cmd, config, cdc, mbm, genAccIterator, outputDir, chainID,
				minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddress, numValidators)
		},
	}

	prepareFlagsForTestnetCmd(cmd)

	return cmd
}

func prepareFlagsForTestnetCmd(cmd *cobra.Command) {
	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "cetd",
		"Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "cetcli",
		"Home directory of the node's cli configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(
		client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("%s%s", authx.DefaultMinGasPriceLimit, dex.DefaultBondDenom), //20sato.CET
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 20cet)")
}

func initTestnet(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	mbm dex.OrderedBasicManager, genAccIterator GenesisAccountsIterator,
	outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome,
	nodeCLIHome, startingIPAddress string, numValidators int) error {

	if chainID == "" {
		chainID = integrationTestChainID + cmn.RandStr(6)
	}

	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]crypto.PubKey, numValidators)
	accs := make([]genaccounts.GenesisAccount, numValidators)
	genFiles := make([]string, numValidators)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeInfo, err := initTestnetNode(cmd, config, cdc,
			outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome,
			nodeCLIHome, startingIPAddress, i)
		if err != nil {
			return err
		}

		nodeIDs[i] = nodeInfo.nodeID
		valPubKeys[i] = nodeInfo.valPubKey
		accs[i] = nodeInfo.acc
		genFiles[i] = nodeInfo.genFile
	}

	if err := initGenFiles(cdc, mbm, chainID, accs, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genAccIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initTestnetNode(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddr string, i int,
) (testnetNodeInfo, error) {

	nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
	nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
	clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
	gentxsDir := filepath.Join(outputDir, "gentxs")

	config.SetRoot(nodeDir)
	config.RPC.ListenAddress = "tcp://0.0.0.0:26657"

	if err := mkNodeHomeDirs(outputDir, nodeDir, clientDir); err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	config.Moniker = nodeDirName
	adjustBlockCommitSpeed(config)

	ip, err := getIP(i, startingIPAddr)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(config)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	memo := fmt.Sprintf("%s@%s:26656", nodeID, ip)
	genFile := config.GenesisFile()

	buf := bufio.NewReader(cmd.InOrStdin())
	prompt := fmt.Sprintf(
		"Password for account '%s' (default %s):", nodeDirName, app.DefaultKeyPass,
	)

	keyPass, err := client.GetPassword(prompt, buf)
	if err != nil && keyPass != "" {
		// An error was returned that either failed to read the password from
		// STDIN or the given password is not empty but failed to meet minimum
		// length requirements.
		return testnetNodeInfo{}, err
	}

	if keyPass == "" {
		keyPass = app.DefaultKeyPass
	}

	addr, secret, err := server.GenerateSaveCoinKey(clientDir, nodeDirName, keyPass, true)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	info := map[string]string{"secret": secret}

	cliPrint, err := json.Marshal(info)
	if err != nil {
		return testnetNodeInfo{}, err
	}

	// save private key seed words
	err = writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint)
	if err != nil {
		return testnetNodeInfo{}, err
	}

	minSelfDel := sdk.NewInt(10000e8)
	accStakingTokens := minSelfDel.MulRaw(10)
	acc := genaccounts.GenesisAccount{
		Address: addr,
		Coins: sdk.Coins{
			sdk.NewCoin(dex.DefaultBondDenom, accStakingTokens),
		},
	}

	msg := staking.NewMsgCreateValidator(
		sdk.ValAddress(addr),
		valPubKey,
		sdk.NewCoin(dex.DefaultBondDenom, minSelfDel),
		staking.NewDescription(nodeDirName, "", "", ""),
		staking.NewCommissionRates(sdk.NewDecWithPrec(3, 2), sdk.OneDec(), sdk.NewDecWithPrec(1, 2)),
		minSelfDel,
	)
	kb, err := keys.NewKeyBaseFromDir(clientDir)
	if err != nil {
		return testnetNodeInfo{}, err
	}
	tx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
	txBldr := auth.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

	signedTx, err := txBldr.SignStdTx(nodeDirName, app.DefaultKeyPass, tx, false)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	txBytes, err := cdc.MarshalJSON(signedTx)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	// gather gentxs folder
	err = writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes)
	if err != nil {
		_ = os.RemoveAll(outputDir)
		return testnetNodeInfo{}, err
	}

	dexConfig := srvconfig.DefaultConfig()
	dexConfig.MinGasPrices = minGasPrices

	configFilePath := filepath.Join(nodeDir, "config/cetd.toml")
	srvconfig.WriteConfigFile(configFilePath, dexConfig)
	return testnetNodeInfo{
		nodeID:    nodeID,
		valPubKey: valPubKey,
		acc:       acc,
		genFile:   genFile,
	}, nil
}

func mkNodeHomeDirs(outputDir, nodeDir, clientDir string) error {
	if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
		_ = os.RemoveAll(outputDir)
		return err
	}

	if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
		_ = os.RemoveAll(outputDir)
		return err
	}

	return nil
}

func initGenFiles(cdc *codec.Codec, mbm dex.OrderedBasicManager, chainID string,
	accs []genaccounts.GenesisAccount, genFiles []string, numValidators int) error {

	appGenState := mbm.DefaultGenesis()

	// set the accounts in the genesis state
	appGenState = genaccounts.SetGenesisStateInAppState(cdc, appGenState, accs)

	addCetTokenForTesting(cdc, appGenState, testnetTokenSupply, accs[0].Address)
	modifyGenStateForTesting(cdc, appGenState, testnetMinSelfDelegation)

	accs = assureTokenDistributionInGenesis(accs, testnetTokenSupply)
	appGenState[genaccounts.ModuleName] = cdc.MustMarshalJSON(accs)

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func assureTokenDistributionInGenesis(accs []genaccounts.GenesisAccount, testnetSupply sdk.Int) []genaccounts.GenesisAccount {
	distributedTokens := sdk.ZeroInt()
	for _, acc := range accs {
		distributedTokens = distributedTokens.Add(acc.Coins[0].Amount)
	}

	if testnetSupply.GT(distributedTokens) {
		accs = append(accs, genaccounts.GenesisAccount{
			Address: sdk.AccAddress(crypto.AddressHash([]byte("left_tokens"))),
			Coins: sdk.Coins{
				sdk.NewCoin(dex.DefaultBondDenom, testnetSupply.Sub(distributedTokens)),
			},
		})
	}
	return accs
}

func modifyGenStateForTesting(cdc *codec.Codec, appGenState map[string]json.RawMessage, testnetMinSelfDelegation int64) {

	var stakingxData stakingx.GenesisState
	cdc.MustUnmarshalJSON(appGenState[stakingx.ModuleName], &stakingxData)

	stakingxData.Params.MinSelfDelegation = testnetMinSelfDelegation
	stakingxData.Params.MinMandatoryCommissionRate = sdk.NewDecWithPrec(2, 2)

	appGenState[stakingx.ModuleName] = cdc.MustMarshalJSON(stakingxData)
}

func addCetTokenForTesting(cdc *codec.Codec,
	appGenState map[string]json.RawMessage, tokenTotalSupply sdk.Int, cetOwner sdk.AccAddress) {

	var assetData asset.GenesisState
	cdc.MustUnmarshalJSON(appGenState[asset.ModuleName], &assetData)

	baseToken, _ := asset.NewToken("CoinEx Chain Native Token",
		dex.CET,
		tokenTotalSupply,
		cetOwner,
		false,
		true,
		false,
		false,
		"www.coinex.org",
		"A public chain built for the decentralized exchange",
		"CF1FAAA36A78BE02",
	)

	var token asset.Token = baseToken
	assetData.Tokens = []asset.Token{token}

	appGenState[asset.ModuleName] = cdc.MustMarshalJSON(assetData)
}

func collectGenFiles(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	nodeIDs []string, valPubKeys []crypto.PubKey,
	numValidators int, outputDir, nodeDirPrefix, nodeDaemonHome string,
	genAccIterator GenesisAccountsIterator) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		moniker := nodeDirName
		config.Moniker = nodeDirName

		config.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutil.NewInitConfig(chainID, gentxsDir, moniker, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr, i)
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = cmn.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}
