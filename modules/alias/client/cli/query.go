package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	"github.com/coinexchain/dex/modules/authx/client/cliutil"
)

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query alias params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func QueryAddressCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-of-alias [alias]",
		Short: "query the corresponding address of an alias",
		Long: `query the corresponding address of an alias. 

Example : 
	cetcli query alias address-of-alias super_super_boy`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
			param := &keepers.QueryAliasInfoParam{Alias: args[0], QueryOp: keepers.GetAddressFromAlias}
			return cliutil.CliQuery(cdc, query, param)
		},
	}
}

func QueryAliasCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "aliases-of-address [address]",
		Short: "query the aliases of an address",
		Long: `query the aliases of an address. 

Example : 
	cetcli query alias aliases-of-address coinex1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
			acc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			param := &keepers.QueryAliasInfoParam{Owner: acc, QueryOp: keepers.ListAliasOfAccount}
			return cliutil.CliQuery(cdc, query, param)
		},
	}
}
