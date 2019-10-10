package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

const FlagAsDefault = "as-default"

func AliasAddCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [alias]",
		Short: "Add an alias for current account",
		Long: `Add an alias for current account.

Example: 
	 cetcli tx alias add super_super_boy --from local_user_1 --as-default yes
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			asDefaultStr := viper.GetString(FlagAsDefault)
			asDefault := true
			if asDefaultStr == "no" {
				asDefault = false
			} else if asDefaultStr != "yes" {
				return errors.New("Invalid value for --as-default, only 'yes' and 'no' are valid")
			}
			msg := &types.MsgAliasUpdate{
				Alias:     args[0],
				IsAdd:     true,
				AsDefault: asDefault,
			}

			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(FlagAsDefault, "yes", "This alias will be used as a default alias or not")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	//cmd.MarkFlagRequired(FlagAsDefault)
	return cmd
}

func AliasRemoveCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [alias]",
		Short: "Remove an alias for current account",
		Long: `Remove an alias for current account.

Example: 
	 cetcli tx alias remove super_super_boy --from local_user_1
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := &types.MsgAliasUpdate{
				Alias: args[0],
				IsAdd: false,
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	return cmd
}
