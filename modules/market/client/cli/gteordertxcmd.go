package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagSymbol    = "symbol"
	FlagOrderType = "order-type"
	FlagPrice     = "price"
	FlagQuantity  = "quantity"
	FlagSide      = "side"
	FlagOrderID   = "orderid"
	FlagUserAddr  = "address"
)

var createOrderFlags = []string{
	FlagSymbol,
	FlagOrderType,
	FlagPrice,
	FlagQuantity,
	FlagSide,
	FlagPricePrecision,
}

func CreateIOCOrderTxCmd(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "createiocorder",
		Short: "",
		Long: `Create an IOC order and sign tx, broadcast to nodes.
		Example:
		$ cetcli tx market creategteoreder --symbol="btc/cet"
		--order-type=2 \
		--price=520 \
		--quantity=10000000 \
		--side=1 \
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadCastOrder(cdc, false)
		},
	}

	markCreateOrderFlags(cmd)
	return cmd
}

func CreateGTEOrderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creategteoreder",
		Short: "Create an GTE order and sign tx",
		Long: `Create an GTE order and sign tx, broadcast to nodes.

Example:
$ cetcli tx market creategteoreder --symbol="btc/cet"
	--order-type=2 \
	--price=520 \
	--quantity=10000000 \
	--side=1 \
	--time-in-force=1000
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAndBroadCastOrder(cdc, true)
		},
	}

	markCreateOrderFlags(cmd)
	return cmd
}

func createAndBroadCastOrder(cdc *codec.Codec, isGTE bool) error {
	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

	sender := cliCtx.GetFromAddress()
	sequence, err := cliCtx.GetAccountSequence(sender)
	if err != nil {
		return err
	}

	msg, err := parseCreateOrderFlags(sender, sequence)
	if err != nil {
		return errors.Errorf("tx flag is error, please see help : " +
			"$ cetcli tx market creategteoreder -h")
	}
	if err = msg.ValidateBasic(); err != nil {
		return err
	}

	symbols := strings.Split(msg.Symbol, market.SymbolSeparator)
	userToken := symbols[0]
	if msg.Side == match.BUY {
		userToken = symbols[1]
	}

	account, err := cliCtx.GetAccount(sender)
	if err != nil {
		return err
	}
	if !account.GetCoins().IsAllGTE(sdk.Coins{sdk.NewCoin(userToken, sdk.NewInt(msg.Quantity))}) {
		return errors.New("No have insufficient cet to create market in blockchain")
	}

	msg.TimeInForce = market.IOC
	if isGTE {
		msg.TimeInForce = market.GTE
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

func parseCreateOrderFlags(sender sdk.AccAddress, sequence uint64) (*market.MsgCreateOrder, error) {
	for _, flag := range createOrderFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, please see help : "+
				"$ cetcli tx market creategteoreder -h", flag)
		}
	}

	msg := &market.MsgCreateOrder{
		Sender:         sender,
		Symbol:         viper.GetString(FlagSymbol),
		OrderType:      byte(viper.GetInt(FlagOrderType)),
		Side:           byte(viper.GetInt(FlagSide)),
		Price:          viper.GetInt64(FlagPrice),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Quantity:       viper.GetInt64(FlagQuantity),
		Sequence:       sequence,
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

	for _, flag := range createOrderFlags {
		cmd.MarkFlagRequired(flag)
	}
}

func QueryOrderCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orderinfo",
		Short: "Query order info",
		Long:  "cetcli query market orderinfo --symbol=[eth/cet] --orderid=[orderid]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			bz, err := cdc.MarshalJSON(market.NewQueryOrderParam(viper.GetString(FlagOrderID)))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.MarketKey, market.QueryOrder)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	markQueryOrDelCmd(cmd)
	return cmd
}

func QueryUserOrderList(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "userorderlist [userAddress]",
		Short: "Query user order list in blockchain",
		Long: "Example:" +
			"cetcli query market userorderlist [userAddress]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			fmt.Println(string(args[0]))
			bz, err := cdc.MarshalJSON(market.QueryUserOrderList{User: args[0]})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.MarketKey, market.QueryUserOrders)
			fmt.Println(route)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

func CancleOrder(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancelorder",
		Short: "cancel order in blockchain",
		Long: "Examples:" +
			"cetcli tx market cancelorder --orderid=[id]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()

			orderid := viper.GetString(FlagOrderID)
			contents := strings.Split(orderid, "-")
			if len(contents) != 2 {
				return errors.Errorf("")
			}

			if bytes.Equal(sender, []byte(contents[0])) {
				return errors.Errorf("")
			}

			if sequence, err := strconv.Atoi(contents[1]); err != nil || sequence < 0 {
				return errors.Errorf("")
			}

			msg := market.MsgCancelOrder{
				Sender:  sender,
				OrderID: orderid,
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	markQueryOrDelCmd(cmd)
	return cmd
}

func markQueryOrDelCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagSymbol, "", "The trading market symbol")
	cmd.Flags().String(FlagOrderID, "", "The order id")
	cmd.MarkFlagRequired(FlagOrderID)
	cmd.MarkFlagRequired(FlagSymbol)
}
