package client

import (
	"github.com/coinexchain/dex/x/asset"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	assCli "github.com/coinexchain/dex/x/asset/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group asset queries under a subcommand
	assQueryCmd := &cobra.Command{
		Use:   asset.ModuleName,
		Short: "Querying commands for the asset module",
	}

	assQueryCmd.AddCommand(client.GetCommands(
		assCli.GetTokenCmd(mc.storeKey, mc.cdc),
	)...)

	return assQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	assTxCmd := &cobra.Command{
		Use:   asset.ModuleName,
		Short: "Asset transactions subcommands",
	}

	assTxCmd.AddCommand(client.PostCommands(
		assCli.IssueTokenCmd(mc.cdc),
	)...)

	return assTxCmd
}
