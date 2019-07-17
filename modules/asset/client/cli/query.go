package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

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
		GetTokenCmd(types.QuerierRoute, cdc),
		GetTokenListCmd(types.QuerierRoute, cdc),
		GetWhitelistCmd(types.QuerierRoute, cdc),
		GetForbiddenAddrCmd(types.QuerierRoute, cdc),
		GetReservedSymbolsCmd(types.QuerierRoute, cdc),
	)...)

	return assQueryCmd
}

// GetTokenCmd returns a query token that will display the info of the
// token at a given token symbol
func GetTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]

			bz, err := cdc.MarshalJSON(types.NewQueryAssetParams(symbol))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetTokenListCmd returns all token that will display
func GetTokenListCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTokenList)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetWhitelistCmd returns whitelist
func GetWhitelistCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]

			bz, err := cdc.MarshalJSON(types.NewQueryWhitelistParams(symbol))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryWhitelist)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetForbiddenAddrCmd returns forbidden addr
func GetForbiddenAddrCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]

			bz, err := cdc.MarshalJSON(types.NewQueryForbiddenAddrParams(symbol))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryForbiddenAddr)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetReservedSymbolsCmd returns reserved symbol list
func GetReservedSymbolsCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryReservedSymbols)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}
