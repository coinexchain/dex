package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount] [unlocked_time]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			time, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bankx.NewMsgSend(from, to, coins, time)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
