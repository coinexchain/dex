package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	authxutils "github.com/coinexchain/dex/modules/authx/client/utils"
	"github.com/coinexchain/dex/modules/distributionx/types"
)

// DonateTxCmd will create a DonateToCommunityPool tx and sign it with the given key.
func DonateTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate [amount]",
		Short: "Donate to community pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc) //.WithAccountDecoder(cdc)

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgDonateToCommunityPool(from, coins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			generateUnsignedTx := viper.GetBool(authxutils.FlagGenerateUnsignedTx)
			if generateUnsignedTx {
				return authxutils.PrintUnsignedTx(cliCtx, txBldr, []sdk.Msg{msg}, from)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.MarkFlagRequired(client.FlagFrom)
	cmd.Flags().Bool(authxutils.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	return cmd
}
