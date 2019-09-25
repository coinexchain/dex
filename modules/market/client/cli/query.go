package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	mktQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands for the market module",
	}
	mktQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(cdc),
		QueryMarketCmd(cdc),
		QueryMarketListCmd(cdc),
		QueryOrderCmd(cdc),
		QueryUserOrderList(cdc))...)
	return mktQueryCmd
}

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func QueryMarketListCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "trading-pairs",
		Short: "query all trading-pair infos in blockchain",
		Long: `query all trading-pair infos in blockchain.

Example :
	cetcli query market trading-pairs
	--trust-node=true --chain-id=coinexdex`,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarkets)
			return cliutil.CliQuery(cdc, query, nil)
		},
	}
}

func QueryMarketCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "trading-pair",
		Short: "query trading-pair info in blockchain",
		Long: `query trading-pair info in blockchain. 

Example : 
	cetcli query market trading-pair 
	eth/cet --trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(strings.Split(args[0], types.SymbolSeparator)) != 2 {
				return errors.Errorf("trading-pair illegal : %s, For example : eth/cet.", args[0])
			}
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
			return cliutil.CliQuery(cdc, query, keepers.NewQueryMarketParam(args[0]))
		},
	}
}

func QueryOrderCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order-info",
		Short: "Query order info in blockchain",
		Long: `Query order info in blockchain. 

Example :
	cetcli query market order-info [orderID] 
	--trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orderID := args[0]
			if len(strings.Split(orderID, types.OrderIDSeparator)) != types.OrderIDPartsNum {
				return fmt.Errorf("order-id is incorrect")
			}
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
			return cliutil.CliQuery(cdc, route, keepers.NewQueryOrderParam(orderID))
		},
	}

	return cmd
}

func QueryUserOrderList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order-list [userAddress]",
		Short: "Query user order list in blockchain",
		Long: `Query user order list in blockchain. 

Example:
	cetcli query market order-list [userAddress] 
	--trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
			return cliutil.CliQuery(cdc, route, keepers.QueryUserOrderList{User: args[0]})
		},
	}

	return cmd
}
