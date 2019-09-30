package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
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
