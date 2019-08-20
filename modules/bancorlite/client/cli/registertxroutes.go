package cli

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	bancorliteQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands to get the information of a symbol pair's bancor pool",
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
