package cli

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

const FlagAsDefault = "as-default"

func AliasAddCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an alias for current account",
		Long: `Add an alias for current account.

Example: 
	 cetcli tx alias add --as-default yes super_super_boy 
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()
			asDefaultStr := viper.GetString(FlagAsDefault)
			asDefault := true
			if asDefaultStr == "no" {
				asDefault = false
			} else if asDefaultStr != "yes" {
				return errors.New("Invalid value for --as-default, only 'yes' and 'no' are valid")
			}
			msg := &types.MsgAliasUpdate{
				Owner:     sender,
				Alias:     args[0],
				IsAdd:     true,
				AsDefault: asDefault,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagAsDefault, "yes", "This alias will be used as a default alias or not")
	cmd.MarkFlagRequired(FlagAsDefault)
	return cmd
}

func AliasRemoveCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an alias for current account",
		Long: `Remove an alias for current account.

Example: 
	 cetcli tx alias remove super_super_boy 
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()
			msg := &types.MsgAliasUpdate{
				Owner: sender,
				Alias: args[0],
				IsAdd: false,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
