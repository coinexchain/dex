package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetTokenCmd returns a query token that will display the info of the
// token at a given token symbol
// nolint: unparam
func GetTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-info [symbol]",
		Short: "Query token info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			symbol := args[0]

			bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(symbol))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryToken)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var token asset.Token
			cdc.MustUnmarshalJSON(res, &token)
			return cliCtx.PrintOutput(token)
		},
	}
	return cmd
}
