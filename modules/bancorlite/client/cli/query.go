package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query bancorlite params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func QueryBancorInfoCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info [stock] [money]",
		Short: "query the banor pool's information about a symbol pair",
		Long: `query the banor pool's information about a symbol pair. 

Example : 
	cetcli query bancorlite info stock money --trust-node=true --chain-id=coinexdex`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryBancorInfo)
			symbol := dex.GetSymbol(args[0], args[1])
			param := &keepers.QueryBancorInfoParam{Symbol: symbol}
			return cliutil.CliQuery(cdc, query, param)
		},
	}
}
