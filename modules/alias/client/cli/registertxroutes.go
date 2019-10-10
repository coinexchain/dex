package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
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
