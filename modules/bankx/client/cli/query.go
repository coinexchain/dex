package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/go-amino"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	aliasQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands for the bankx module",
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
		Short: "Query bankx params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keeper.QueryParameters)
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
