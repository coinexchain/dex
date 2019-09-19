package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx/internal/types"
)

// get the root query command of this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	assQueryCmd := &cobra.Command{
		Use:   auth.ModuleName,
		Short: "Querying commands for the auth module",
	}

	assQueryCmd.AddCommand(client.GetCommands(
		GetQueryParamsCmd(cdc),
	)...)

	return assQueryCmd
}

func GetQueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query auth params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.MergedParams
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

func GetAccountXCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			accRetriever := auth.NewAccountRetriever(cliCtx)
			if err = accRetriever.EnsureExists(key); err != nil {
				return err
			}

			acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key)
			if err != nil {
				return err
			}

			aux, err := GetAccountX(cliCtx, key)
			if err != nil { // it's ok
				aux = types.AccountX{}
			}

			all := types.AccountAll{Account: acc, AccountX: aux}

			return cliCtx.PrintOutput(all)
		},
	}
	return client.GetCommands(cmd)[0]
}

func GetAccountX(ctx context.CLIContext, address []byte) (types.AccountX, error) {
	res, err := queryAccountX(ctx, address)
	if err != nil {
		return types.AccountX{}, err
	}

	var accountX types.AccountX
	if err := ctx.Codec.UnmarshalJSON(res, &accountX); err != nil {
		return types.AccountX{}, err
	}

	return accountX, nil
}

func queryAccountX(ctx context.CLIContext, addr sdk.AccAddress) ([]byte, error) {
	bz, err := ctx.Codec.MarshalJSON(auth.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAccountX)

	res, _, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}
