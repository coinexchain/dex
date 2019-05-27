package cli

import (
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
)

func CreateGTEOrderTxCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creategteoreder",
		Short: "",
		Long:  "",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			accout, err := cliCtx.GetAccount([]byte(args[0]))
			if err != nil {
				return errors.Errorf("Not find account with the addr : %s", args[0])
			}

			sequence, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			symbols := strings.Split(args[2], market.SymbolSeparator)
			userToken := symbols[0]
			side, err := strconv.Atoi(args[7])
			if err != nil || (side != match.BUY && side != match.SELL) {
				return errors.New("side out of range")
			}
			if side == match.BUY {
				userToken = symbols[1]
			}

			quantity, err := strconv.Atoi(args[6])
			if err != nil {
				return err
			}
			if !accout.GetCoins().IsAllGTE(sdk.Coins{sdk.NewCoin(userToken, sdk.NewInt(int64(quantity)))}) {
				return errors.New("No have insufficient cet to create market in blockchain")
			}

			ordertype, err := strconv.Atoi(args[3])
			if err != nil || (ordertype != int(market.LimitOrder)) {
				return errors.Errorf("order type out of range")
			}

			pricePrecision, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			price, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}

			timeInForce, err := strconv.Atoi(args[8])
			if err != nil {
				return err
			}

			msg := market.MsgCreateGTEOrder{
				Sender:         []byte(args[0]),
				Sequence:       uint64(sequence),
				Symbol:         args[2],
				OrderType:      byte(ordertype),
				PricePrecision: byte(pricePrecision),
				Price:          int64(price),
				Quantity:       int64(quantity),
				Side:           byte(side),
				TimeInForce:    timeInForce,
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	return cmd
}
