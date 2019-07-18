package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/coinexchain/dex/app"
	"github.com/cosmos/cosmos-sdk/client"
)

func TestInitConfig(t *testing.T) {
	cdc := app.MakeCodec()
	rootCmd := createRootCmd(cdc)
	viper.Set("trust-node", true)
	_ = rootCmd.PersistentFlags().String(cli.HomeFlag, "./", "")

	err := initConfig(rootCmd)
	require.Equal(t, nil, err)
}

func TestFixDescriptions(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "cmd",
		Short: "Command",
	}
	cmd.Flags().String(client.FlagFees, "", "Fees")
	cmd.Flags().String(client.FlagGasPrices, "", "Fees")
	//fixDescriptions(cmd)
	//require.Equal(t, cmd.Flag(client.FlagFees).Usage, "Fees to pay along with transaction; eg: 100cet")
}
