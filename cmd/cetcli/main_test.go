package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"

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
