package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/staking"
	staking_cli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	stakingQueryCmd := staking_cli.GetQueryCmd(staking.QuerierRoute, cdc)
	bondPoolCmd := client.GetCommands(GetCmdQueryPool(cdc))[0]
	paramsCmd := client.GetCommands(GetCmdQueryParams(cdc))[0]

	//replace pool cmd with new bondPoolCmd which can also show the non-bondable-cet-tokens in locked positions
	cmd := replacePoolCmd(stakingQueryCmd, "pool", bondPoolCmd)
	cmd = replacePoolCmd(stakingQueryCmd, "params", paramsCmd)
	return cmd
}

func replacePoolCmd(stakingQueryCmd *cobra.Command, use string, BondPoolCmd *cobra.Command) *cobra.Command {
	var oldPoolCmd *cobra.Command
	for _, cmd := range stakingQueryCmd.Commands() {
		if cmd.Use == use {
			oldPoolCmd = cmd
		}
	}

	stakingQueryCmd.RemoveCommand(oldPoolCmd)
	stakingQueryCmd.AddCommand(BondPoolCmd)

	return stakingQueryCmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current staking pool values",
		Long: strings.TrimSpace(`Query values for amounts stored in the staking pool:

$ cetcli query staking pool
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData("custom/stakingx/pool", nil)
			if err != nil {
				return err
			}

			println(string(res))
			return nil
		},
	}
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current staking parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as staking parameters.

Example:
$ %s query staking params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params, err := queryStakingParams(cdc, cliCtx)
			if err != nil {
				return err
			}

			paramsx, err := queryStakingXParams(cdc, cliCtx)
			if err != nil {
				return err
			}

			mergedParams := MergedParams{
				UnbondingTime:              params.UnbondingTime,
				MaxValidators:              params.MaxValidators,
				MaxEntries:                 params.MaxEntries,
				BondDenom:                  params.BondDenom,
				MinSelfDelegation:          paramsx.MinSelfDelegation,
				MinMandatoryCommissionRate: paramsx.MinMandatoryCommissionRate,
			}
			return cliCtx.PrintOutput(mergedParams)
		},
	}
}

func queryStakingParams(cdc *codec.Codec, cliCtx context.CLIContext) (staking.Params, error) {
	route := fmt.Sprintf("custom/%s/%s", staking.StoreKey, staking.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return staking.Params{}, err
	}
	var params staking.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

func queryStakingXParams(cdc *codec.Codec, cliCtx context.CLIContext) (types.Params, error) {
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, staking.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}
