package cli

import (
	"fmt"
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

var (
	FlagSymbol      = "symbol"
	FlagOrderType   = "order-type"
	FlagPrice       = "price"
	FlagQuantity    = "quantity"
	FlagSide        = "side"
	FlagTimeInForce = "time-in-force"
)

var createGTEOrderFlags = []string{
	FlagSymbol,
	FlagOrderType,
	FlagPrice,
	FlagQuantity,
	FlagSide,
	FlagTimeInForce,
}

func CreateGTEOrderTxCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			sender := cliCtx.GetFromAddress()
			sequence, err := cliCtx.GetAccountSequence(sender)
			if err != nil {
				return err
			}

			msg, err := parseCreateOrderFlags(sender, sequence)
			if err != nil {
				return errors.Errorf("tx flag is error, pls see help : " +
					"$ cetcli tx market creategteoreder -h")
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			symbols := strings.Split(args[2], market.SymbolSeparator)
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

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "The trading market symbol")
	cmd.Flags().Int(FlagOrderType, -1, "The order type limited to 2")
	cmd.Flags().Int(FlagPrice, -1, "The price in the order")
	cmd.Flags().Int(FlagQuantity, -1, "The number of tokens will be trade in the order ")
	cmd.Flags().Int(FlagSide, -1, "The side in the order")
	cmd.Flags().Int(FlagTimeInForce, -1, "The time-in-force for GTE order")

	for _, flag := range createGTEOrderFlags {
		cmd.MarkFlagRequired(flag)
	}

	return cmd
}

func parseCreateOrderFlags(sender sdk.AccAddress, sequence uint64) (*market.MsgCreateGTEOrder, error) {
	for _, flag := range createGTEOrderFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, pls see help : "+
				"$ cetcli tx market creategteoreder -h", flag)
		}
	}

	msg := &market.MsgCreateGTEOrder{
		Sender:         sender,
		Symbol:         viper.GetString(FlagSymbol),
		OrderType:      byte(viper.GetInt(FlagOrderType)),
		Side:           byte(viper.GetInt(FlagSide)),
		Price:          viper.GetInt64(FlagPrice),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Quantity:       viper.GetInt64(FlagQuantity),
		TimeInForce:    viper.GetInt(FlagTimeInForce),
		Sequence:       sequence,
	}

	return msg, nil
}
