package cli

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

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
