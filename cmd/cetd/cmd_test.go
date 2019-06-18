package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/cmd"
)

func TestCreateRootCmd(t *testing.T) {
	cmd.InitSdkConfig()
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := createRootCmd(ctx, cdc)
	require.Equal(t, 7, len(rootCmd.Commands()))
}
