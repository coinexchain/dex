package client

import (
	"github.com/coinexchain/dex/modules/stakingx/client/cli"
	"github.com/coinexchain/dex/modules/stakingx/client/rest"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	staking_cli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	staking_types "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/types"
)

func NewStakingXModuleClient() types.ModuleClient {
	return StakingXModuleClient{}
}

type StakingXModuleClient struct {
}

func (StakingXModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (StakingXModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return staking_cli.GetTxCmd(staking_types.StoreKey, cdc)
}

func (StakingXModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	stakingQueryCmd := staking_cli.GetQueryCmd(staking_types.QuerierRoute, cdc)
	BondPoolCmd := client.GetCommands(cli.GetCmdQueryPool(cdc))[0]

	//replace pool cmd with new bondPoolCmd which can also show the non-bondable-cet-tokens in locked positions
	return replacePoolCmd(stakingQueryCmd, BondPoolCmd)
}

func replacePoolCmd(stakingQueryCmd *cobra.Command, BondPoolCmd *cobra.Command) *cobra.Command {
	var oldPoolCmd *cobra.Command
	for _, cmd := range stakingQueryCmd.Commands() {
		if cmd.Use == "pool" {
			oldPoolCmd = cmd
		}
	}

	stakingQueryCmd.RemoveCommand(oldPoolCmd)
	stakingQueryCmd.AddCommand(BondPoolCmd)

	return stakingQueryCmd
}
