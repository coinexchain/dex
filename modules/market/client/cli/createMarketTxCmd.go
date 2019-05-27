package cli

import (
	"fmt"
	"strconv"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	FlagCreator        = "creator"
	FlagStock          = "stock"
	FlagMonet          = "money"
	FlagPricePrecision = "price-precision"
)

func CreateMarketCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createmarket ",
		Short: "generate tx to create market",
		Long: "generate a tx and sign it to create market in dex blockchain." +
			"Example : createmarket [creator] [stock] [money] [priceprecision]",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			accout, err := cliCtx.GetAccount([]byte(args[0]))
			if err != nil {
				return errors.Errorf("Not find account with the addr : %s", args[0])
			}

			if !accout.GetCoins().IsAllGTE(sdk.Coins{market.CreateMarketSpendCet}) {
				return errors.New("No have insufficient cet to create market in blockchain")
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryTokenList)
			res, _ := cliCtx.QueryWithData(route, nil)
			if res == nil {
				return errors.New("Not query asset info from blockchain")
			}

			var tokens []asset.Token
			cdc.MustUnmarshalJSON(res, &tokens)

			if !IsExistStockAndMoneySymbol(args[1], args[1], tokens) {
				return errors.New("stock or monry is not exist in blockchain")
			}

			pricePrecision, err := strconv.Atoi(args[3])
			if err != nil || (pricePrecision < market.MinimumTokenPricePrecision ||
				pricePrecision > market.MaxTokenPricePrecision) {
				return errors.Errorf("price precision out of range [%d, %d]",
					market.MinimumTokenPricePrecision, market.MaxTokenPricePrecision)
			}

			msg := market.NewMsgCreateMarketInfo(args[1], args[2], []byte(args[0]), byte(pricePrecision))

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagCreator, "", "The address to create a trading market")
	cmd.Flags().String(FlagStock, "", "The exist token symbol as stock")
	cmd.Flags().String(FlagMonet, "", "")

	return cmd
}

func IsExistStockAndMoneySymbol(stock, money string, tokens []asset.Token) bool {
	var (
		findStock bool
		findMoney bool
	)

	for _, t := range tokens {
		if stock == t.GetSymbol() {
			findStock = true
		} else if money == t.GetSymbol() {
			findMoney = true
		}
	}

	if findMoney && findStock {
		return true
	}
	return false
}
