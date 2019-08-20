package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	aliasQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands for the alias module",
	}
	aliasQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(cdc),
		QueryAliasCmd(cdc),
		QueryAddressCmd(cdc),
	)...)
	return aliasQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
	aliasTxCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "alias transactions subcommands",
	}

	aliasTxCmd.AddCommand(client.PostCommands(
		AliasAddCmd(cdc),
		AliasRemoveCmd(cdc),
	)...)

	return aliasTxCmd
}
