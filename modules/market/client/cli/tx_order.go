package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

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
	 cetcli tx market create-ioc-order --trading-pair=btc/cet 
	--order-type=2 --price=520 --quantity=10000000 
	--side=1 --price-precision=10 --from=bob --identify=1
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
	cetcli tx market create-gte-order --trading-pair=btc/cet 
	--order-type=2 --price=520 --quantity=10000000 --side=1 
	--price-precision=10 --blocks=<100000> --from=bob --identify=1
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
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

	sender := cliCtx.GetFromAddress()
	msg, err := parseCreateOrderFlags(sender)
	if err != nil {
		if isGTE {
			return errors.Errorf("tx flag is error, please see help : " +
				"$ cetcli tx market create-gte-order -h")
		}
		return errors.Errorf("tx flag is error, please see help : " +
			"$ cetcli tx market create-ioc-order -h")
	}
	msg.TimeInForce = types.IOC
	if isGTE {
		msg.TimeInForce = types.GTE
	}
	if err = msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

func parseCreateOrderFlags(sender sdk.AccAddress) (*types.MsgCreateOrder, error) {
	for _, flag := range createOrderFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop" + flag)
		}
	}
	blocks := types.DefaultGTEOrderLifetime
	if viper.GetInt(FlagBlocks) > 0 {
		blocks = viper.GetInt(FlagBlocks)
	}

	msg := &types.MsgCreateOrder{
		Sender:         sender,
		Identify:       byte(viper.GetInt(FlagIdentify)),
		TradingPair:    viper.GetString(FlagSymbol),
		OrderType:      byte(viper.GetInt(FlagOrderType)),
		Side:           byte(viper.GetInt(FlagSide)),
		Price:          viper.GetInt64(FlagPrice),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Quantity:       viper.GetInt64(FlagQuantity),
		ExistBlocks:    blocks,
	}

	return msg, nil
}

func markCreateOrderFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagSymbol, "", "The trading market symbol")
	cmd.Flags().Int(FlagOrderType, -1, "The order type limited to 2")
	cmd.Flags().Int(FlagPrice, -1, "The price in the order")
	cmd.Flags().Int(FlagQuantity, -1, "The number of tokens will be trade in the order ")
	cmd.Flags().Int(FlagSide, -1, "The side in the order")
	cmd.Flags().Int(FlagPricePrecision, -1, "The price precision in order")
	cmd.Flags().Int(FlagIdentify, 0, "Because a transaction can contain multiple order "+
		"creation messages, the identify field was added to the order creation message to give each "+
		"order a unique ID. So the order ID consists of user address, user sequence, identify,")

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

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			sender := cliCtx.GetFromAddress()
			orderid := viper.GetString(FlagOrderID)
			msg, err := CheckSenderAndOrderID(sender, orderid)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	markQueryOrDelCmd(cmd)
	return cmd
}

func CheckSenderAndOrderID(sender []byte, orderID string) (types.MsgCancelOrder, error) {
	var (
		addr sdk.AccAddress
		err  error
		msg  types.MsgCancelOrder
	)

	contents := strings.Split(orderID, types.OrderIDSeparator)
	if len(contents) != types.OrderIDPartsNum {
		return msg, errors.Errorf(" illegal order-id")
	}

	if addr, err = sdk.AccAddressFromBech32(contents[0]); err != nil {
		return msg, err
	}
	if !bytes.Equal(addr, sender) {
		return msg, errors.Errorf("sender address is not match order sender, sender : %s, order issuer : %s", sender, addr)
	}

	if sequence, err := strconv.Atoi(contents[1]); err != nil || sequence < 0 {
		return msg, errors.Errorf("illegal order sequence, actual %d", sequence)
	}

	msg = types.MsgCancelOrder{
		Sender:  sender,
		OrderID: orderID,
	}

	return msg, nil
}

func markQueryOrDelCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagOrderID, "", "The order id")
	cmd.MarkFlagRequired(FlagOrderID)
}
