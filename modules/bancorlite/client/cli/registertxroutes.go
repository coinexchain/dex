package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	bancorliteQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands for the bancorlite module",
	}
	bancorliteQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(cdc),
		QueryBancorInfoCmd(cdc),
	)...)
	return bancorliteQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
	bancorliteTxCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "bancorlite transactions subcommands",
	}

	bancorliteTxCmd.AddCommand(client.PostCommands(
		BancorInitCmd(cdc),
		BancorTradeCmd(cdc),
		BancorCancelCmd(cdc),
	)...)

	return bancorliteTxCmd
}
