package cli

import (
	"fmt"
	"github.com/coinexchain/dex/client/cliutil"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	aliasQueryCmd := &cobra.Command{
		Use:   bank.ModuleName,
		Short: "Querying commands for the bank module",
	}
	aliasQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(cdc),
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
