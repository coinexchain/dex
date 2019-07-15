package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func TestCreateRootCmd(t *testing.T) {
	dex.InitSdkConfig()

	rootCmd := createCetdCmd()
	require.Equal(t, 15, len(rootCmd.Commands()))
}

func TestNewApp(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	cet := newApp(logger, db, log.NewSyncWriter(os.Stdout))
	value := reflect.ValueOf(cet).Interface().(*app.CetChainApp)
	require.Equal(t, "CetChainApp", value.Name())
}
