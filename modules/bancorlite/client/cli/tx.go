package cli

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

const (
	FlagMaxSupply          = "max-supply"
	FlagMaxMoney           = "max-money"
	FlagStockPrecision     = "stock-precision"
	FlagMaxPrice           = "max-price"
	FlagSide               = "side"
	FlagAmount             = "amount"
	FlagMoneyLimit         = "money-limit"
	FlagInitPrice          = "init-price"
	FlagEarliestCancelTime = "earliest-cancel-time"
)

var bancorInitFlags = []string{
	FlagMaxSupply,
	FlagMaxMoney,
	FlagStockPrecision,
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
	 cetcli tx bancorlite init stock money --max-supply=10000000000000 --max-money=100000 --stock-precision=3 --max-price=5 --init-price=1 --earliest-cancel-time=1563954165
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			maxPrice, err0 := sdk.NewDecFromStr(viper.GetString(FlagMaxPrice))
			if err0 != nil || maxPrice.IsZero() {
				return errors.New("max price is invalid or zero")
			}
			initPrice, err0 := sdk.NewDecFromStr(viper.GetString(FlagInitPrice))
			if err0 != nil || initPrice.IsNegative() {
				return errors.New("init price is negative")
			}
			maxSupply, ok := sdk.NewIntFromString(viper.GetString(FlagMaxSupply))
			if !ok {
				return errors.New("max supply is invalid")
			}
			maxMoney, ok := sdk.NewIntFromString(viper.GetString(FlagMaxMoney))
			if !ok {
				return errors.New("max money is invalid")
			}
			precision, convertErr := strconv.Atoi(viper.GetString(FlagStockPrecision))
			if convertErr != nil {
				return errors.New("stock precision is invalid")
			}
			time, err := strconv.ParseInt(viper.GetString(FlagEarliestCancelTime), 10, 64)
			if err != nil {
				return errors.New("bancor earliest-cancel-time is invalid")
			}
			msg := &types.MsgBancorInit{
				Stock:              args[0],
				Money:              args[1],
				InitPrice:          viper.GetString(FlagInitPrice),
				MaxSupply:          maxSupply,
				StockPrecision:     byte(precision),
				MaxPrice:           viper.GetString(FlagMaxPrice),
				MaxMoney:           maxMoney,
				EarliestCancelTime: time,
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(FlagMaxSupply, "0", "The maximum supply of this pool.")
	cmd.Flags().String(FlagMaxMoney, "0", "The maximum money of this pool")
	cmd.Flags().String(FlagStockPrecision, "0", "The precision of stock")
	cmd.Flags().String(FlagMaxPrice, "0", "The maximum reachable price when all the supply are sold out")
	cmd.Flags().String(FlagEarliestCancelTime, "0", "The time that bancor can be canceled")
	cmd.Flags().String(FlagInitPrice, "0", "The init price of this bancor")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")
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
				Stock:      args[0],
				Money:      args[1],
				Amount:     viper.GetInt64(FlagAmount),
				IsBuy:      isBuy,
				MoneyLimit: viper.GetInt64(FlagMoneyLimit),
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().Int(FlagAmount, 0, "The amount of tokens to be traded.")
	cmd.Flags().Int(FlagMoneyLimit, 0, "The upper bound of money you want to pay when buying, or the lower bound of money you want to get when selling. Specify zero or negative value if you do not want a such a limit.")
	cmd.Flags().String(FlagSide, "", "the side of the trade, 'buy' or 'sell'.")
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

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
			msg := &types.MsgBancorCancel{
				Stock: args[0],
				Money: args[1],
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}
	cmd.Flags().Bool(cliutil.FlagGenerateUnsignedTx, false, "Generate a unsigned tx")

	return cmd
}
