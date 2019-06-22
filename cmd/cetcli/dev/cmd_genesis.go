package dev

import (
	"fmt"

	"github.com/spf13/cobra"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/app"
)

func ExampleGenesisCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-genesis",
		Short: "Print Cetd example genesis JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			genState := createExampleGenesisState(cdc)
			return printGenesisState(cdc, genState, "coinexdex-1")
		},
	}
	return cmd
}

func TestnetGenesisCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet-genesis",
		Short: "Print Cetd testnet genesis JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			genState := createTestnetGenesisState(cdc)
			return printGenesisState(cdc, genState, "coinexdex-test1")
		},
	}
	return cmd
}

func printGenesisState(cdc *codec.Codec, genState app.GenesisState, chainID string) error {
	gneStateBytes, err := codec.MarshalJSONIndent(cdc, genState)
	if err != nil {
		return err
	}

	genDoc := tm.GenesisDoc{
		ChainID:    chainID,
		Validators: nil,
		AppState:   gneStateBytes,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	genDoc.ConsensusParams.Evidence.MaxAge = app.DefaultEvidenceMaxAge

	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(genDocBytes))
	return nil
}
