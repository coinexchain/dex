package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
)

func init() {
	dex.InitSdkConfig()
}

func TestInitConfig(t *testing.T) {
	cdc := app.MakeCodec()
	rootCmd := createRootCmd(cdc)
	viper.Set("trust-node", true)
	_ = rootCmd.PersistentFlags().String(cli.HomeFlag, "./", "")

	err := initConfig(rootCmd)
	require.Equal(t, nil, err)
}

func TestFixDescriptions(t *testing.T) {
	cdc := app.MakeCodec()
	rootCmd := createRootCmd(cdc)
	fixDescriptions(rootCmd)
	// TODO
	//require.Equal(t, cmd.Flag(client.FlagFees).Usage, "Fees to pay along with transaction; eg: 100cet")
}

// https://trello.com/c/c5sx305m
func TestKeysParseCmd(t *testing.T) {
	tmpFileName := "TestKeysParseCmd.output"
	tmpOutput, err := os.Create(tmpFileName)
	require.NoError(t, err)

	oldStdout := os.Stdout
	os.Stdout = tmpOutput
	defer func() {
		os.Stdout = oldStdout
		os.Remove(tmpFileName)
	}()

	parseCmd := getKeysParseCmd()
	require.NotNil(t, parseCmd)
	viper.Set(cli.OutputFlag, "text")
	err = parseCmd.RunE(nil, []string{"BD7F0C07B8EC2DA56E2BB89C958925B8B13B272C"})
	require.NoError(t, err)
	output, err := ioutil.ReadFile(tmpFileName)
	require.NoError(t, err)
	require.Contains(t, string(output), "coinex")
	require.NotContains(t, string(output), "cosmos")
}

func getKeysParseCmd() *cobra.Command {
	cdc := app.MakeCodec()
	for _, subCmd := range createRootCmd(cdc).Commands() {
		if subCmd.Use == "keys" {
			for _, subCmd := range keys.Commands().Commands() {
				if subCmd.Use == "parse <hex-or-bech32-address>" {
					return subCmd
				}
			}
		}
	}
	return nil
}
