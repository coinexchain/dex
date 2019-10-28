package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

const (
	FlagSymbol    = "trading-pair"
	FlagOrderType = "order-type"
	FlagPrice     = "price"
	FlagQuantity  = "quantity"
	FlagSide      = "side"
	FlagOrderID   = "order-id"
	FlagBlocks    = "blocks"
	FlagTime      = "time"
	FlagIdentify  = "identify"
)

var createOrderFlags = []string{
	FlagSymbol,
	FlagOrderType,
	FlagPrice,
	FlagQuantity,
	FlagSide,
	FlagPricePrecision,
	FlagIdentify,
}

func CreateIOCOrderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-ioc-order",
		Short: "Create an IOC order and sign tx",
		Long: `Create an IOC order and sign tx, broadcast to nodes.

Example: 
	 cetcli tx market create-ioc-order --trading-pair=btc/cet \
	--order-type=2 --price=520 --quantity=10000000 \
	--side=1 --price-precision=10 --from=bob --identify=1 \
	--chain-id=coinexdex --gas=10000 --fees=1000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadCastOrder(cdc, false)
		},
	}
	markCreateOrderFlags(cmd)
	return cmd
}

func CreateGTEOrderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gte-order",
		Short: "Create an GTE order and sign tx",
		Long: `Create an GTE order and sign tx, broadcast to nodes. 

Example:
	cetcli tx market create-gte-order --trading-pair=btc/cet \
	--order-type=2 --price=520 --quantity=10000000 --side=1 \
	--price-precision=10 --blocks=100000 --from=bob --identify=1 \
	--chain-id=coinexdex --gas=10000 --fees=1000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadCastOrder(cdc, true)
		},
	}
	markCreateOrderFlags(cmd)
	cmd.Flags().Int(FlagBlocks, 10000, "the gte order will exist at least blocks in blockChain")
	return cmd
}

func createAndBroadCastOrder(cdc *codec.Codec, isGTE bool) error {
	msg, err := parseCreateOrderFlags(isGTE)
	if err != nil {
		if isGTE {
			return errors.Errorf("errors : %s, please see help : "+
				"$ cetcli tx market create-gte-order -h", err.Error())
		}
		return errors.Errorf("errors : %s, please see help : "+
			"$ cetcli tx market create-ioc-order -h", err.Error())
	}
	return cliutil.CliRunCommand(cdc, msg)
}

func parseCreateOrderFlags(isGTE bool) (*types.MsgCreateOrder, error) {
	for _, flag := range createOrderFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop" + flag)
		}
	}
	msg := &types.MsgCreateOrder{
		Identify:       byte(viper.GetInt(FlagIdentify)),
		TradingPair:    viper.GetString(FlagSymbol),
		OrderType:      byte(viper.GetInt(FlagOrderType)),
		Side:           byte(viper.GetInt(FlagSide)),
		Price:          viper.GetInt64(FlagPrice),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Quantity:       viper.GetInt64(FlagQuantity),
		ExistBlocks:    int64(viper.GetInt(FlagBlocks)),
		TimeInForce:    types.IOC,
	}
	if isGTE {
		msg.TimeInForce = types.GTE
	}
	return msg, nil
}

func markCreateOrderFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagSymbol, "", "The trading pair symbol")
	cmd.Flags().Int(FlagOrderType, 2, "The identify of the price limit : 2; (Currently, only price limit orders are supported)")
	cmd.Flags().Int(FlagPrice, 100, "The price of the order")
	cmd.Flags().Int(FlagQuantity, 100, "The number of tokens will be trade in the order ")
	cmd.Flags().Int(FlagSide, 1, "The buying or selling direction of an order.(buy : 1; sell : 2)")
	cmd.Flags().Int(FlagPricePrecision, 8, "The price precision in the order")
	cmd.Flags().Int(FlagIdentify, 0, "A transaction can contain multiple order "+
		"creation messages, the identify field was added to the order creation message to give each "+
		"order a unique ID. So the order ID consists of user address, user sequence, identify.")

	for _, flag := range createOrderFlags {
		cmd.MarkFlagRequired(flag)
	}
}

func CancelOrder(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-order",
		Short: "cancel order in blockchain",
		Long: `cancel order in blockchain. 

Examples:
	cetcli tx market cancel-order --order-id=[id] 
	--trust-node=true --from=bob --chain-id=coinexdex`,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := &types.MsgCancelOrder{
				OrderID: viper.GetString(FlagOrderID),
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}
	markQueryOrDelCmd(cmd)
	return cmd
}

func markQueryOrDelCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagOrderID, "", "The order id")
	cmd.MarkFlagRequired(FlagOrderID)
}
