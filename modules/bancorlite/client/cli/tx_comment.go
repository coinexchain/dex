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

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

const (
	FlagMaxSupply  = "max-supply"
	FlagMaxPrice   = "max-price"
	FlagSide       = "side"
	FlagAmount     = "amount"
	FlagMoneyLimit = "money-limit"
)

var bancorInitFlags = []string{
	FlagMaxSupply,
	FlagMaxPrice,
}

var bancorTradeFlags = []string{
	FlagSide,
	FlagAmount,
	FlagMoneyLimit,
}

func BancorInitCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a bancor pool for a token",
		Long: `Initialize a bancor pool for a token, specifying the maximum supply of this pool and the maximum reachable price when all the supply are sold out.

Example: 
	 cetcli tx bancorlite init cetdac --max-supply=10000000000000 --max-price=5 
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()
			var maxPrice sdk.Dec
			types.FillDec(viper.GetString(FlagMaxPrice), &maxPrice)
			if maxPrice.IsZero() {
				return errors.New("Max Price is Invalid or Zero")
			}
			maxSupply, ok := sdk.NewIntFromString(viper.GetString(FlagMaxSupply))
			if !ok {
				return errors.New("Max Supply is Invalid")
			}
			msg := &types.MsgBancorInit{
				Owner:     sender,
				Token:     args[0],
				MaxSupply: maxSupply,
				MaxPrice:  maxPrice,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagMaxSupply, "0", "The maximum supply of this pool.")
	cmd.Flags().String(FlagMaxPrice, "0", "the maximum reachable price when all the supply are sold out")

	for _, flag := range bancorInitFlags {
		cmd.MarkFlagRequired(flag)
	}
	return cmd
}

func BancorTradeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trade",
		Short: "Trade with a bancor pool",
		Long: `Sell tokens to a bancor pool or buy tokens from a bancor pool.

Example: 
	 cetcli tx bancorlite trade cetdac --side buy --amount=100 --money-limit=120
	 cetcli tx bancorlite trade cetdac --side sell --amount=100 --money-limit=80
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()
			var isBuy bool
			switch viper.GetString(FlagSide) {
			case "buy":
				isBuy = true
			case "sell":
				isBuy = false
			default:
				return errors.New("Unknown Side. Please specify 'buy' or 'sell'")
			}
			msg := &types.MsgBancorTrade{
				Sender:     sender,
				Token:      args[0],
				Amount:     viper.GetInt64(FlagAmount),
				IsBuy:      isBuy,
				MoneyLimit: viper.GetInt64(FlagMoneyLimit),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagAmount, 0, "The amount of tokens to be traded.")
	cmd.Flags().Int(FlagMoneyLimit, 0, "The upper bound of money you want to pay when buying, or the lower bound of money you want to get when selling. Specify zero or negative value if you do not want a such a limit.")
	cmd.Flags().String(FlagSide, "", "the side of the trade, 'buy' or 'sell'.")

	for _, flag := range bancorTradeFlags {
		cmd.MarkFlagRequired(flag)
	}
	return cmd
}
