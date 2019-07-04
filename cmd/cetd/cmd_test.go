package main

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func TestCreateRootCmd(t *testing.T) {
	dex.InitSdkConfig()
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := createCetdCmd(ctx, cdc)
	require.Equal(t, 15, len(rootCmd.Commands()))
}

func TestNewApp(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	cet := newApp(logger, db, log.NewSyncWriter(os.Stdout))
	value := reflect.ValueOf(cet).Interface().(*app.CetChainApp)
	require.Equal(t, "CetChainApp", value.Name())
}
