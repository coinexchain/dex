package client

import (
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/client/cli"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

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
		Use:   market.MarketKey,
		Short: "Querying commands for the market module",
	}
	assQueryCmd.AddCommand(client.GetCommands(
		cli.QueryMarketCmd(mc.cdc),
		cli.QueryOrderCmd(mc.cdc))...)
	return assQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	assTxCmd := &cobra.Command{
		Use:   market.MarketKey,
		Short: "market transactions subcommands",
	}

	assTxCmd.AddCommand(client.PostCommands(
		cli.CreateMarketCmd(mc.storeKey, mc.cdc),
		cli.CreateGTEOrderTxCmd(mc.storeKey, mc.cdc),
	)...)

	return assTxCmd
}
