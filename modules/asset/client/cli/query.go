package cli

import (
	"fmt"
	"strings"

	"github.com/coinexchain/dex/modules/asset"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

// get the root query command of this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	assQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the asset module",
	}

	assQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(types.QuerierRoute, cdc),
		GetCmdQueryToken(types.QuerierRoute, cdc),
		GetCmdQueryTokenList(types.QuerierRoute, cdc),
		GetCmdQueryTokenWhitelist(types.QuerierRoute, cdc),
		GetCmdQueryTokenForbiddenAddr(types.QuerierRoute, cdc),
		GetCmdQueryTokenReservedSymbols(types.QuerierRoute, cdc),
	)...)

	return assQueryCmd
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query asset params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

// GetCmdQueryToken returns a query token that will display the info of the
// token at a given token symbol
func GetCmdQueryToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [symbol]",
		Short: "Query token info",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a token. You can find the token by token symbol".

Example:
$ cetcli query asset token abc
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryToken)
			symbol := args[0]
			if err := asset.ValidateTokenSymbol(symbol); err != nil {
				return err
			}
			params := types.NewQueryAssetParams(symbol)
			return cliutil.CliQuery(cdc, route, params)
		},
	}
	return cmd
}

// GetCmdQueryTokenList returns all token that will display
func GetCmdQueryTokenList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens ",
		Short: "Query all token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all tokens. 
Example:
$ cetcli query asset tokens
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryTokenList)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
	return cmd
}

// GetCmdQueryTokenWhitelist returns whitelist
func GetCmdQueryTokenWhitelist(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist [symbol]",
		Short: "Query whitelist",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query whitelist for a token. You can find it by token symbol".

Example:
$ cetcli query asset whitelist abc
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryWhitelist)
			symbol := args[0]
			if err := asset.ValidateTokenSymbol(symbol); err != nil {
				return err
			}
			params := types.NewQueryWhitelistParams(symbol)
			return cliutil.CliQuery(cdc, route, params)
		},
	}
	return cmd
}

// GetCmdQueryTokenForbiddenAddr returns forbidden addr
func GetCmdQueryTokenForbiddenAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forbidden-addresses [symbol]",
		Short: "Query forbidden addresses",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query forbidden addresses for a token. You can find it by token symbol".

Example:
$ cetcli query asset forbidden-addresses abc
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryForbiddenAddr)
			symbol := args[0]
			if err := asset.ValidateTokenSymbol(symbol); err != nil {
				return err
			}
			params := types.NewQueryForbiddenAddrParams(symbol)
			return cliutil.CliQuery(cdc, route, params)
		},
	}
	return cmd
}

// GetCmdQueryTokenReservedSymbols returns reserved symbol list
func GetCmdQueryTokenReservedSymbols(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserved-symbols",
		Short: "Query reserved symbols",
		Args:  cobra.ExactArgs(0),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query reserved symbols list".

Example:
$ cetcli query asset reserved-symbols
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryReservedSymbols)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
	return cmd
}
