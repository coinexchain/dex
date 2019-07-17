package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/comment/internal/types"
	"github.com/coinexchain/dex/modules/comment/internal/keepers"
)

func QueryCommentCountCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-count",
		Short: "query total comment count in blockchain",
		Long: `query total comment count in blockchain. 

Example : 
	cetcli query comment get-count --trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)//.WithAccountDecoder(cdc)
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryCommentCount)
			res, _, err := cliCtx.QueryWithData(query, nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
