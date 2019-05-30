package cli

import (
	clientx "github.com/coinexchain/dex/client"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
)

func GetAccountXCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).WithAccountDecoder(cdc)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if err = cliCtx.EnsureAccountExistsFromAddr(key); err != nil {
				return err
			}

			acc, err := cliCtx.GetAccount(key)
			if err != nil {
				return err
			}

			aux, err := clientx.GetAccountX(cliCtx, key)
			if err != nil {
				return err
			}

			all := authx.AccountAll{Account: acc, AccountX: aux}

			return cliCtx.PrintOutput(all)
		},
	}
	return client.GetCommands(cmd)[0]
}
