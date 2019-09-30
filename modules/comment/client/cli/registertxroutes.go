package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/modules/comment/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	commentTxCmd := &cobra.Command{
		Use:   types.StoreKey,
		Short: "comment transactions subcommands",
	}

	commentTxCmd.AddCommand(client.PostCommands(
		CreateNewThreadCmd(cdc),
		CreateFollowupCommentCmd(cdc),
		RewardCommentsCmd(cdc),
	)...)

	return commentTxCmd
}
