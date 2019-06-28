package client

import (
	"github.com/coinexchain/dex/modules/stakingx/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	stakingclient "github.com/cosmos/cosmos-sdk/x/staking/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	stakingMC *stakingclient.ModuleClient
	storeKey  string
	cdc       *amino.Codec
}

func NewModuleClient(mc *stakingclient.ModuleClient, storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{mc, storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	stakingQueryCmd := mc.stakingMC.GetQueryCmd()
	BondPoolCmd := client.GetCommands(cli.GetCmdQueryPool(mc.storeKey, mc.cdc))[0]

	//replace pool cmd with new bondPoolCmd which can also show the non-bondable-cet-tokens in locked positions
	return mc.replacePoolCmd(stakingQueryCmd, BondPoolCmd)
}

func (mc ModuleClient) replacePoolCmd(stakingQueryCmd *cobra.Command, BondPoolCmd *cobra.Command) *cobra.Command {
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

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	return mc.stakingMC.GetTxCmd()
}
