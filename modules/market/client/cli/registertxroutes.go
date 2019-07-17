package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	mktQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying commands for the market module",
	}
	mktQueryCmd.AddCommand(client.GetCommands(
		QueryMarketCmd(cdc),
		QueryOrderCmd(cdc),
		QueryUserOrderList(cdc))...)
	// cli.QueryWaitCancelMarkets(mc.cdc))...)
	return mktQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
	mktTxCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "market transactions subcommands",
	}

	mktTxCmd.AddCommand(client.PostCommands(
		CreateMarketCmd(cdc),
		CreateGTEOrderTxCmd(cdc),
		CreateIOCOrderTxCmd(cdc),
		CancelOrder(cdc),
		CancelMarket(cdc),
		ModifyTradingPairPricePrecision(cdc),
	)...)

	return mktTxCmd
}
