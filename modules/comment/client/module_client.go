package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/coinexchain/dex/modules/comment"
	"github.com/coinexchain/dex/modules/comment/client/cli"
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
	commentQueryCmd := &cobra.Command{
		Use:   comment.StoreKey,
		Short: "Querying command to get the total comment count",
	}
	commentQueryCmd.AddCommand(client.GetCommands(
		cli.QueryCommentCountCmd(mc.cdc))...)
	return commentQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	commentTxCmd := &cobra.Command{
		Use:   comment.StoreKey,
		Short: "comment transactions subcommands",
	}

	commentTxCmd.AddCommand(client.PostCommands(
		cli.RewardCommentsCmd(mc.cdc),
		cli.CreateNewThreadCmd(mc.cdc),
		cli.CreateFollowupCommentCmd(mc.cdc),
	)...)

	return commentTxCmd
}
