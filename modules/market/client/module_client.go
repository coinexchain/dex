package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/client/cli"
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
	mktQueryCmd := &cobra.Command{
		Use:   market.StoreKey,
		Short: "Querying commands for the market module",
	}
	mktQueryCmd.AddCommand(client.GetCommands(
		cli.QueryMarketCmd(mc.cdc),
		cli.QueryOrderCmd(mc.cdc),
		cli.QueryUserOrderList(mc.cdc))...)
	// cli.QueryWaitCancelMarkets(mc.cdc))...)
	return mktQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	mktTxCmd := &cobra.Command{
		Use:   market.StoreKey,
		Short: "market transactions subcommands",
	}

	mktTxCmd.AddCommand(client.PostCommands(
		cli.CreateMarketCmd(mc.storeKey, mc.cdc),
		cli.CreateGTEOrderTxCmd(mc.cdc),
		cli.CreateIOCOrderTxCmd(mc.cdc),
		cli.CancelOrder(mc.cdc),
		cli.CancelMarket(mc.cdc),
		cli.ModifyTradingPairPricePrecision(mc.cdc),
	)...)

	return mktTxCmd
}
