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
	FlagUnlockTime  = "unlock-time"
	FlagToAddress   = "to-address"
	FlagFromAddress = "from-address"
	FlagSupervisor  = "supervisor"
	FlagReward      = "reward"
	FlagOperation   = "operation"
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
    cetcli tx send supervised-tx 1000000000cet \
        --from-address=coinex1hckjvduhckfaxq2tuythfd270cex94c0hv5hs7 \
        --to-address=coinex1ke3qq22zvzlcdh3j8nenlrjxmvnrna7z426n0x \
        --supervisor=coinex12agppqdn8tr3wme40uxex2fnxf7780j7f3zpxn \
        --unlock-time=1600000000 \
        --reward=100000000 \
        --operation=1 \
        --from=local_user
`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse coins trying to be sent
			coin, err := sdk.ParseCoin(args[0])
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

			from := viper.GetString(FlagFromAddress)
			fromAddr, err := sdk.AccAddressFromBech32(from)
			if err != nil {
				return err
			}

			to := viper.GetString(FlagToAddress)
			toAddr, err := sdk.AccAddressFromBech32(to)
			if err != nil {
				return err
			}

			supervisor := viper.GetString(FlagSupervisor)
			supervisorAddr, err := sdk.AccAddressFromBech32(supervisor)
			if err != nil {
				return err
			}

			reward := viper.GetInt64(FlagReward)
			if reward < 0 || coin.Amount.LT(sdk.NewInt(reward)) {
				return fmt.Errorf("invalid reward: %d", reward)
			}

			operation := byte(viper.GetInt(FlagOperation))
			if operation < types.Create || operation > types.EarlierUnlockBySupervisor {
				return fmt.Errorf("invalid operation type: %d", operation)
			}

			msg := types.NewMsgSupervisedSend(fromAddr, supervisorAddr, toAddr, coin, unlockTime, reward, operation)
			return cliutil.CliRunCommand(cdc, &msg)
		},
	}

	cmd.Flags().Int64(FlagUnlockTime, 0, "The unix timestamp when tokens can transfer again")
	cmd.Flags().String(FlagFromAddress, "", "The address which the locked amount comes from")
	cmd.Flags().String(FlagToAddress, "", "The address which the locked amount is sent to")
	cmd.Flags().String(FlagSupervisor, "", "The supervisor's address")
	cmd.Flags().Int64(FlagReward, 0, "The reward for supervisor")
	cmd.Flags().Int(FlagOperation, 0, "Operation type (create : 0; return : 1; unlock by sender: 2; unlock by supervisor: 3)")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	_ = cmd.MarkFlagRequired(FlagUnlockTime)
	_ = cmd.MarkFlagRequired(FlagFromAddress)
	_ = cmd.MarkFlagRequired(FlagToAddress)
	_ = cmd.MarkFlagRequired(FlagSupervisor)
	_ = cmd.MarkFlagRequired(FlagReward)
	_ = cmd.MarkFlagRequired(FlagOperation)

	return cmd
}
