package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/authx/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
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
			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func GetAccountXCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountMix)
			acc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			param := auth.NewQueryAccountParams(acc)
			return cliutil.CliQuery(cdc, route, &param)
		},
	}
}
