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
	FlagSender     = "sender"
	FlagSupervisor = "supervisor"
	FlagReward     = "reward"
	FlagOperation  = "operation"
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

	cmd.AddCommand(client.PostCommands(
		SendSupervisedTxCmd(cdc),
	)...)

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

// SendSupervisedTxCmd
func SendSupervisedTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supervised-tx [amount]",
		Short: "Create and sign a supervised tx",
		Long: `Create and sign a supervised tx.

Example:
    cetcli tx send supervised-tx coinex1ke3qq22zvzlcdh3j8nenlrjxmvnrna7z426n0x 1000000000cet \
        --sender=coinex1hckjvduhckfaxq2tuythfd270cex94c0hv5hs7 \
        --unlock-time=1600000000 \
        --reward=100000000 \
        --operation=1 \
        --from=local_user
`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			unlockTime := viper.GetInt64(FlagUnlockTime)
			currentTime := time.Now().Unix()
			if unlockTime < currentTime {
				return fmt.Errorf("unlock time should be later than the current time")
			}

			var fromAddr sdk.AccAddress
			sender := viper.GetString(FlagSender)
			if sender != "" {
				if fromAddr, err = sdk.AccAddressFromBech32(sender); err != nil {
					return err
				}
			}

			var supervisorAddr sdk.AccAddress
			supervisor := viper.GetString(FlagSupervisor)
			if supervisor != "" {
				if supervisorAddr, err = sdk.AccAddressFromBech32(supervisor); err != nil {
					return err
				}
			}

			reward := viper.GetInt64(FlagReward)
			operation := byte(viper.GetInt(FlagOperation))

			msg := types.NewMsgSupervisedSend(fromAddr, supervisorAddr, toAddr, coin, unlockTime, reward, operation)
			return cliutil.CliRunCommand(cdc, &msg)
		},
	}

	cmd.Flags().Int64(FlagUnlockTime, 0, "The unix timestamp when tokens can transfer again")
	cmd.Flags().String(FlagSender, "", "The supervised amount sender's address (required when return or unlock-by-supervisor)")
	cmd.Flags().String(FlagSupervisor, "", "The supervisor's address (required when create or unlock-by-sender if there is a supervisor)")
	cmd.Flags().Int64(FlagReward, 0, "The reward for supervisor")
	cmd.Flags().Int(FlagOperation, 0, "Operation type (create: 0; return: 1; unlock-by-sender: 2; unlock-by-supervisor: 3)")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	_ = cmd.MarkFlagRequired(FlagUnlockTime)
	_ = cmd.MarkFlagRequired(FlagOperation)

	return cmd
}
