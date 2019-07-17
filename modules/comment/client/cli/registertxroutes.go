package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/coinexchain/dex/modules/comment/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	commentQueryCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "Querying command to get the total comment count",
	}
	commentQueryCmd.AddCommand(client.GetCommands(
		QueryCommentCountCmd(cdc))...)
	return commentQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
	commentTxCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "comment transactions subcommands",
	}

	commentTxCmd.AddCommand(client.PostCommands(
		RewardCommentsCmd(cdc),
		CreateNewThreadCmd(cdc),
		CreateFollowupCommentCmd(cdc),
	)...)

	return commentTxCmd
}
