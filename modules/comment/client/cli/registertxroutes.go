package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/coinexchain/dex/modules/comment/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *amino.Codec) *cobra.Command {
	return nil
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *amino.Codec) *cobra.Command {
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
