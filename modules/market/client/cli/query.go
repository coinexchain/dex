package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market"
)

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
			cliCtx := context.NewCLIContext().WithCodec(cdc)//.WithAccountDecoder(cdc)
			if len(strings.Split(args[0], market.SymbolSeparator)) != 2 {
				return errors.Errorf("symbol illegal : %s, For example : eth/cet.", args[0])
			}

			bz, err := cdc.MarshalJSON(market.NewQueryMarketParam(args[0]))
			if err != nil {
				return err
			}
			query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryMarket)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)//.WithAccountDecoder(cdc)

			time, err := strconv.Atoi(args[0])
			if time <= 0 || err != nil {
				return errors.Errorf("Invalid unix time")
			}

			bz, err := cdc.MarshalJSON(market.QueryCancelMarkets{Time: int64(time)})
			if err != nil {
				return err
			}

			query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryWaitCancelMarkets)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)//.WithAccountDecoder(cdc)

			orderID := args[0]
			if len(strings.Split(orderID, market.OrderIDSeparator)) != 3 {
				return fmt.Errorf("order-id is incorrect")
			}

			bz, err := cdc.MarshalJSON(market.NewQueryOrderParam(orderID))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryOrder)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	markQueryOrDelCmd(cmd)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)//.WithAccountDecoder(cdc)

			queryAddr := args[0]
			if _, err := sdk.AccAddressFromBech32(queryAddr); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(market.QueryUserOrderList{User: queryAddr})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryUserOrders)
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
