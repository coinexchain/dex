package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

func QueryBancorInfoCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "query the banor pool's information about a symbol pair",
		Long: `query the banor pool's information about a symbol pair. 

Example : 
	cetcli query bancorlite info stock money --trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryBancorInfo)
			symbol := args[0] + "/" + args[1]
			param := &keepers.QueryBancorInfoParam{Symbol: symbol}

			bz, err := cdc.MarshalJSON(param)

			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
