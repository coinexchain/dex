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
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/match"
)

const (
	FlagSymbol    = "symbol"
	FlagOrderType = "order-type"
	FlagPrice     = "price"
	FlagQuantity  = "quantity"
	FlagSide      = "side"
	FlagOrderID   = "order-id"
	FlagUserAddr  = "address"
	FlagHeight    = "height"
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
		Use:   "create-ioc-order",
		Short: "Create an IOC order and sign tx",
		Long: "Create an IOC order and sign tx, broadcast to nodes. \n" +
			"Example:" +
			"$ cetcli tx market create-ioc-order --symbol=btc/cet " +
			"--order-type=2 --price=520 --quantity=10000000 " +
			"--side=1 --price-precision=10 --from=bob " +
			"--chain-id=coinexdex --gas=10000 --fees=1000cet",
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
		Long: "Create an GTE order and sign tx, broadcast to nodes. \n" +
			"Example:" +
			"$ cetcli tx market create-gte-order --symbol=btc/cet " +
			"--order-type=2 --price=520 --quantity=10000000 --side=1 " +
			"--price-precision=10 --from=bob --chain-id=coinexdex " +
			"--gas=10000 --fees=1000cet",
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
		if isGTE {
			return errors.Errorf("tx flag is error, please see help : " +
				"$ cetcli tx market create-gte-order -h")
		}
		return errors.Errorf("tx flag is error, please see help : " +
			"$ cetcli tx market create-ioc-order -h")
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
			return nil, fmt.Errorf("--%s flag is a noop" + flag)
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
		Use:   "order-info",
		Short: "Query order info in blockchain",
		Long: "Query order info in blockchain. \n" +
			"Example : " +
			"cetcli query market order-info " +
			"--order-id=[orderID] --trust-node=true --chain-id=coinexdex",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			bz, err := cdc.MarshalJSON(market.NewQueryOrderParam(viper.GetString(FlagOrderID)))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryOrder)
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
		Use:   "user-order-list [userAddress]",
		Short: "Query user order list in blockchain",
		Long: "Query user order list in blockchain. \n" +
			"Example:" +
			"cetcli query market user-order-list --address=[userAddress] --trust-node=true --chain-id=coinexdex",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			queryAddr := viper.GetString(FlagUserAddr)
			if _, err := sdk.AccAddressFromBech32(queryAddr); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(market.QueryUserOrderList{User: queryAddr})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryUserOrders)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagUserAddr, "", "The address of the user to be queried")
	cmd.MarkFlagRequired(FlagUserAddr)
	return cmd
}

func CancelOrder(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-order",
		Short: "cancel order in blockchain",
		Long: "cancel order in blockchain. \n" +
			"Examples:" +
			"cetcli tx market cancel-order --order-id=[id] " +
			"--trust-node=true --from=bob --chain-id=coinexdex",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			sender := cliCtx.GetFromAddress()
			orderid := viper.GetString(FlagOrderID)
			msg, err := CheckSenderAndOrderID(sender, orderid)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	markQueryOrDelCmd(cmd)
	return cmd
}

func CheckSenderAndOrderID(sender []byte, orderID string) (market.MsgCancelOrder, error) {
	var (
		addr sdk.AccAddress
		err  error
		msg  market.MsgCancelOrder
	)

	contents := strings.Split(orderID, "-")
	if len(contents) != 2 {
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

	msg = market.MsgCancelOrder{
		Sender:  sender,
		OrderID: orderID,
	}

	return msg, nil
}

func markQueryOrDelCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagOrderID, "", "The order id")
	cmd.MarkFlagRequired(FlagOrderID)
}
