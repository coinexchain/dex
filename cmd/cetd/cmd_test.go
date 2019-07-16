package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
	dbm "github.com/tendermint/tendermint/libs/db"
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
	cet := newApp(logger, db, log.NewSyncWriter(os.Stdout))
	value := reflect.ValueOf(cet).Interface().(*app.CetChainApp)
	require.Equal(t, "CetChainApp", value.Name())
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
