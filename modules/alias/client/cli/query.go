package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

func QueryAddressCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-of-alias",
		Short: "query the corresponding address of an alias",
		Long: `query the corresponding address of an alias. 

Example : 
	cetcli query alias address-of-alias super_super_boy`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
			param := &keepers.QueryAliasInfoParam{Alias: args[0], QueryOp: keepers.GetAddressFromAlias}
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

func QueryAliasCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "aliases-of-address",
		Short: "query the aliases of an address",
		Long: `query the aliases of an address. 

Example : 
	cetcli query alias aliases-of-address coinex1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
			acc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			param := &keepers.QueryAliasInfoParam{Owner: acc, QueryOp: keepers.ListAliasOfAccount}
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
