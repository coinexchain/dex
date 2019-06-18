package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
)

func TestCreateRootCmd(t *testing.T) {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	dex.InitSdkConfig()

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	rootCmd := createRootCmd(cdc)
	require.Equal(t, 13, len(rootCmd.Commands()))
}
