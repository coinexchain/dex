package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	// cli.QueryWaitCancelMarkets(mc.cdc))...)
	return mktQueryCmd
}

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarkets)
			res, _, err := cliCtx.QueryWithData(query, nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
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
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)
			if len(strings.Split(args[0], types.SymbolSeparator)) != 2 {
				return errors.Errorf("symbol illegal : %s, For example : eth/cet.", args[0])
			}

			bz, err := cdc.MarshalJSON(keepers.NewQueryMarketParam(args[0]))
			if err != nil {
				return err
			}
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
			res, _, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

func QueryWaitCancelMarkets(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait-cancel-trading-pair",
		Short: "Query wait cancel trading-pair info in special time",
		Long: `Query wait cancel trading-pair info in special time.

Example:
	cetcli query market 
	wait-cancel-trading-pair --time=10000 
	--trust-node=true --chain-id=coinexdex`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			time, err := strconv.Atoi(args[0])
			if time <= 0 || err != nil {
				return errors.Errorf("Invalid unix time")
			}

			bz, err := cdc.MarshalJSON(keepers.QueryCancelMarkets{Time: int64(time)})
			if err != nil {
				return err
			}

			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryWaitCancelMarkets)
			res, _, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}

			var markets []string
			if err := cdc.UnmarshalJSON(res, &markets); err != nil {
				return err
			}
			fmt.Println(markets)

			return nil
		},
	}

	return cmd
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
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			orderID := args[0]
			if len(strings.Split(orderID, types.OrderIDSeparator)) != types.OrderIDPartsNum {
				return fmt.Errorf("order-id is incorrect")
			}

			bz, err := cdc.MarshalJSON(keepers.NewQueryOrderParam(orderID))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			queryAddr := args[0]
			if _, err := sdk.AccAddressFromBech32(queryAddr); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(keepers.QueryUserOrderList{User: queryAddr})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
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
