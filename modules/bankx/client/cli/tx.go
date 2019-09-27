package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

const (
	FlagUnlockTime = "unlock-time"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			unlockTime := viper.GetInt64(FlagUnlockTime)
			if unlockTime < 0 {
				return fmt.Errorf("invalid unlock time: %d", unlockTime)
			}

			currentTime := time.Now().Unix()
			if unlockTime > 0 && unlockTime < currentTime {
				return fmt.Errorf("unlock time should be later than the current time")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSend(nil, to, coins, unlockTime)
			return cliutil.CliRunCommand(cdc, &msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	_ = cmd.MarkFlagRequired(client.FlagFrom)
	cmd.Flags().Int64(FlagUnlockTime, 0, "The unix timestamp when tokens can transfer again")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")
	return cmd
}

func RequireMemoCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "require-memo <bool>",
		Short: "Mark if memo is required to receive coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			required, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgSetTransferMemoRequired(nil, required)
			return cliutil.CliRunCommand(cdc, &msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	_ = cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
