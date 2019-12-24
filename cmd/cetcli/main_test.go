package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"

	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
)

func init() {
	dex.InitSdkConfig()
}

func TestInitConfig(t *testing.T) {
	rootCmd := newRootCmd()
	viper.Set("trust-node", true)
	_ = rootCmd.PersistentFlags().String(cli.HomeFlag, "./", "")

	err := initConfig(rootCmd)
	require.Equal(t, nil, err)
}

func TestFixDescriptions(t *testing.T) {
	rootCmd := newRootCmd()
	fixDescriptions(rootCmd)
	// TODO
	//require.Equal(t, cmd.Flag(client.FlagFees).Usage, "Fees to pay along with transaction; eg: 100cet")
}

// https://trello.com/c/c5sx305m
func TestKeysParseCmd(t *testing.T) {
	viper.Set(cli.OutputFlag, "text")
	parseCmd := getSubCmd(t, newRootCmd(), "keys", "parse")
	output := execAndGetOutput(t, parseCmd, []string{"BD7F0C07B8EC2DA56E2BB89C958925B8B13B272C"})
	require.Contains(t, output, "coinex")
	require.NotContains(t, output, "cosmos")
}

func TestGovCmd(t *testing.T) {
	rootCmd := newRootCmd()
	cmd := getSubCmd(t, rootCmd, "tx", "gov", "submit-proposal")
	require.NotContains(t, cmd.Long, "10test")
	require.Contains(t, cmd.Long, "10cet")

	cmd = getSubCmd(t, rootCmd, "tx", "gov", "submit-proposal", "param-change")
	require.NotContains(t, cmd.Long, `"stake"`)
	require.Contains(t, cmd.Long, `"cet"`)

	cmd = getSubCmd(t, rootCmd, "tx", "gov", "submit-proposal", "community-pool-spend")
	require.NotContains(t, cmd.Long, `"stake"`)
	require.Contains(t, cmd.Long, `"cet"`)
	require.NotContains(t, cmd.Long, `Atoms`)
	require.Contains(t, cmd.Long, `CETs`)
}

func newRootCmd() *cobra.Command {
	cdc := app.MakeCodec()
	return createRootCmd(cdc)
}

func getSubCmd(t *testing.T,
	cmd *cobra.Command, names ...string) *cobra.Command {

	name := names[0]
	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == name {
			if len(names) == 1 {
				return subCmd
			}
			return getSubCmd(t, subCmd, names[1:]...)
		}
	}
	require.Fail(t, "command not found: "+strings.Join(names, " "))
	return nil
}

func execAndGetOutput(t *testing.T,
	cmd *cobra.Command, args []string) string {

	tmpFileName := "TestKeysParseCmd.output"
	tmpOutput, err := os.Create(tmpFileName)
	require.NoError(t, err)

	oldStdout := os.Stdout
	os.Stdout = tmpOutput
	defer func() {
		os.Stdout = oldStdout
		os.Remove(tmpFileName)
	}()

	err = cmd.RunE(cmd, args)
	require.NoError(t, err)
	output, err := ioutil.ReadFile(tmpFileName)
	require.NoError(t, err)
	return string(output)
}
