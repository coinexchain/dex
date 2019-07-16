package dev

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/spf13/viper"

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

func TestnetGenesisCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet-genesis",
		Short: "Print Cetd testnet genesis JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateGenesisJSON(cdc)
		},
	}

	addCmdFlags(cmd)
	return cmd
}

func addCmdFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagAddrCirculation, "", "circulationn account address")
	cmd.Flags().String(flagAddrCoinExFoundation, "", "coinex foundation account address")
	cmd.Flags().String(flagAddrVesting2020, "", "coinex team vesting account address unfreezed on 2020")
	cmd.Flags().String(flagAddrVesting2021, "", "coinex team vesting account address unfreezed on 2021")
	cmd.Flags().String(flagAddrVesting2022, "", "coinex team vesting account address unfreezed on 2022")
	cmd.Flags().String(flagAddrVesting2023, "", "coinex team vesting account address unfreezed on 2023")
	cmd.Flags().String(flagAddrVesting2024, "", "coinex team vesting account address unfreezed on 2024")
	_ = cmd.MarkFlagRequired(flagAddrCirculation)
	_ = cmd.MarkFlagRequired(flagAddrCoinExFoundation)
	_ = cmd.MarkFlagRequired(flagAddrVesting2020)
	_ = cmd.MarkFlagRequired(flagAddrVesting2021)
	_ = cmd.MarkFlagRequired(flagAddrVesting2022)
	_ = cmd.MarkFlagRequired(flagAddrVesting2023)
	_ = cmd.MarkFlagRequired(flagAddrVesting2024)

	_ = cmd.MarkFlagRequired(client.FlagChainID)
}

func generateGenesisJSON(cdc *codec.Codec) error {
	genState := createTestnetGenesisState(cdc)

	chainID := viper.GetString(client.FlagChainID)
	return printGenesisState(cdc, genState, chainID)
}

func printGenesisState(cdc *codec.Codec, genState map[string]json.RawMessage, chainID string) error {
	orderedGenState := app.NewOrderedGenesisState(genState)
	gneStateBytes, err := codec.MarshalJSONIndent(cdc, orderedGenState)
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
