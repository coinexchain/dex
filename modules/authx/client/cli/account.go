package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/coinexchain/dex/modules/authx"
)

func GetAccountXCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc) //.WithAccountDecoder(cdc)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			accRetriever := authtypes.NewAccountRetriever(cliCtx)
			if err = accRetriever.EnsureExists(key); err != nil {
				return err
			}

			acc, err := authtypes.NewAccountRetriever(cliCtx).GetAccount(key)
			if err != nil {
				return err
			}

			aux, err := GetAccountX(cliCtx, key)
			if err != nil { // it's ok
				aux = authx.AccountX{}
			}

			all := authx.AccountAll{Account: acc, AccountX: aux}

			return cliCtx.PrintOutput(all)
		},
	}
	return client.GetCommands(cmd)[0]
}
