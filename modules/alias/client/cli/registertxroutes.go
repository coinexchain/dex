package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
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
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
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
