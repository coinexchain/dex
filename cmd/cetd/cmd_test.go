package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/server"

	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
)

func init() {
	dex.InitSdkConfig()
}

func TestCreateRootCmd(t *testing.T) {
	rootCmd := createCetdCmd()
	require.Equal(t, 16, len(rootCmd.Commands()))
}

func TestNewApp(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	viper.Set(server.FlagMinGasPrices, "20.0cet")
	cet := newApp(logger, db, log.NewSyncWriter(os.Stdout))
	value := reflect.ValueOf(cet).Interface().(*app.CetChainApp)
	require.Equal(t, "CoinExChainApp", value.Name())
}

func TestInitCmd(t *testing.T) {
	testHmoe := "./testhome"
	defer os.RemoveAll(testHmoe)

	rootCmd := createCetdCmd()
	executor := cli.PrepareBaseCmd(rootCmd, "GA", "here")

	os.Args = []string{"cetd", "init", "mynode"}
	viper.Set("home", testHmoe)
	err := executor.Execute()
	require.Nil(t, err)
}

func TestVersionCmd(t *testing.T) {
	testHmoe := "./testhome"
	_, err := os.Stat(testHmoe)
	require.True(t, os.IsNotExist(err))

	rootCmd := createCetdCmd()
	executor := cli.PrepareBaseCmd(rootCmd, "GA", "here")

	os.Args = []string{"cetd", "version"}
	viper.Set("home", testHmoe)
	err = executor.Execute()
	require.Nil(t, err)

	_, err = os.Stat(testHmoe)
	require.True(t, os.IsNotExist(err))
}
