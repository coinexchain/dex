package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	aliasQueryCmd := &cobra.Command{
		Use:   bank.ModuleName,
		Short: "Querying commands for the bank module",
	}
	aliasQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(cdc),
		QueryBalancesCmd(cdc),
	)...)
	return aliasQueryCmd
}

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query bank params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keeper.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func QueryBalancesCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "balances [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keeper.QueryBalances)
			acc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			param := keeper.NewQueryAddrBalances(acc)
			return cliutil.CliQuery(cdc, route, &param)
		},
	}
}
