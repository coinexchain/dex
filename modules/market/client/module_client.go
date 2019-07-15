package client

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/client/cli"
	"github.com/coinexchain/dex/modules/market/client/rest"
)

// AssetModuleClient exports all client functionality from this module
type ModuleClient struct {
}

func NewModuleClient() ModuleClient {
	return ModuleClient{}
}

func (mc ModuleClient) RegisterRESTRoutes(cliCtx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(cliCtx, rtr, market.ModuleCdc)
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	mktQueryCmd := &cobra.Command{
		Use:   market.StoreKey,
		Short: "Querying commands for the market module",
	}
	mktQueryCmd.AddCommand(client.GetCommands(
		cli.QueryMarketCmd(cdc),
		cli.QueryOrderCmd(cdc),
		cli.QueryUserOrderList(cdc))...)
	// cli.QueryWaitCancelMarkets(mc.cdc))...)
	return mktQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd(cdc *amino.Codec) *cobra.Command {
	mktTxCmd := &cobra.Command{
		Use:   market.StoreKey,
		Short: "market transactions subcommands",
	}

	mktTxCmd.AddCommand(client.PostCommands(
		cli.CreateMarketCmd(cdc),
		cli.CreateGTEOrderTxCmd(cdc),
		cli.CreateIOCOrderTxCmd(cdc),
		cli.CancelOrder(cdc),
		cli.CancelMarket(cdc),
		cli.ModifyTradingPairPricePrecision(cdc),
	)...)

	return mktTxCmd
}
