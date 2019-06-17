package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
)

func TestCreateRootCmd(t *testing.T) {
	initSdkConfig()
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := createRootCmd(ctx, cdc)
	require.Equal(t, 7, len(rootCmd.Commands()))
}
