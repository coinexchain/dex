package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/app"
)

func TestCreateRootCmd(t *testing.T) {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	initSdkConfig()

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	rootCmd := createRootCmd(cdc)
	require.Equal(t, 1, len(rootCmd.Commands()))
}
