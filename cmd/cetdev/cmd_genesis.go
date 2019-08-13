package main

import (
	"fmt"

	"github.com/spf13/cobra"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/app"
)

const (
	flagAddrCirculation      = "addr-circulation"
	flagAddrCoinExFoundation = "addr-coinex-foundation"
	flagAddrVesting2020      = "addr-vesting-2020"
	flagAddrVesting2021      = "addr-vesting-2021"
	flagAddrVesting2022      = "addr-vesting-2022"
	flagAddrVesting2023      = "addr-vesting-2023"
	flagAddrVesting2024      = "addr-vesting-2024"
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
