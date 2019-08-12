package cli

import (
	"errors"
	"strconv"

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
	FlagMaxSupply          = "max-supply"
	FlagMaxPrice           = "max-price"
	FlagSide               = "side"
	FlagAmount             = "amount"
	FlagMoneyLimit         = "money-limit"
	FlagInitPrice          = "init-price"
	FlagEarliestCancelTime = "earliest-cancel-time"
)

var bancorInitFlags = []string{
	FlagMaxSupply,
	FlagMaxPrice,
	FlagEarliestCancelTime,
	FlagInitPrice,
}

var bancorTradeFlags = []string{
	FlagSide,
	FlagAmount,
	FlagMoneyLimit,
}

func BancorInitCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [stock] [money]",
		Short: "Initialize a bancor pool for a stock/money pair",
		Long: `Initialize a bancor pool for a stock/money pair, specifying the maximum supply of this pool and the maximum reachable price when all the supply are sold out, specifying the init price, and specifying the time before which no cancellation is allowed.

Example: 
	 cetcli tx bancorlite init stock money --max-supply=10000000000000 --max-price=5 --init-price=1 --earliest-cancel-time=1563954165
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()
			maxPrice, err0 := sdk.NewDecFromStr(viper.GetString(FlagMaxPrice))
			if err0 != nil || maxPrice.IsZero() {
				return errors.New("max Price is Invalid or Zero")
			}
			initPrice, err0 := sdk.NewDecFromStr(viper.GetString(FlagInitPrice))
			if err0 != nil || initPrice.IsNegative() {
				return errors.New("init price is negative")
			}
			maxSupply, ok := sdk.NewIntFromString(viper.GetString(FlagMaxSupply))
			if !ok {
				return errors.New("max Supply is Invalid")
			}
			time, err := strconv.ParseInt(viper.GetString(FlagEarliestCancelTime), 10, 64)
			if err != nil {
				return errors.New("bancor enable-cancel-time is invalid")
			}
			msg := &types.MsgBancorInit{
				Owner:              sender,
				Stock:              args[0],
				Money:              args[1],
				InitPrice:          initPrice,
				MaxSupply:          maxSupply,
				MaxPrice:           maxPrice,
				EarliestCancelTime: time,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagMaxSupply, "0", "The maximum supply of this pool.")
	cmd.Flags().String(FlagMaxPrice, "0", "The maximum reachable price when all the supply are sold out")
	cmd.Flags().String(FlagEarliestCancelTime, "0", "The time that bancor can be canceled")
	cmd.Flags().String(FlagInitPrice, "0", "The init price of this bancor")
	for _, flag := range bancorInitFlags {
		cmd.MarkFlagRequired(flag)
	}
	return cmd
}

func BancorTradeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trade [stock] [money]",
		Short: "Trade with a bancor pool",
		Long: `Sell Stocks to a bancor pool or buy Stocks from a bancor pool.

Example: 
	 cetcli tx bancorlite trade stock money --side buy --amount=100 --money-limit=120
	 cetcli tx bancorlite trade stock money --side sell --amount=100 --money-limit=80
`,
		Args: cobra.ExactArgs(2),
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
				return errors.New("unknown Side. Please specify 'buy' or 'sell'")
			}
			msg := &types.MsgBancorTrade{
				Sender:     sender,
				Stock:      args[0],
				Money:      args[1],
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

func BancorCancelCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [stock] [money]",
		Short: "Cancel a bancor pool for a stock/money pair",
		Long: `Cancel a bancor pool for a stock/money pair, sender must be this stock owner

Example: 
	 cetcli tx bancorlite cancel stock money
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sender := cliCtx.GetFromAddress()

			msg := &types.MsgBancorCancel{
				Owner: sender,
				Stock: args[0],
				Money: args[1],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
