package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	FlagStock          = "stock"
	FlagMoney          = "money"
	FlagPricePrecision = "price-precision"
)

var createMarketFlags = []string{
	FlagStock,
	FlagMoney,
	FlagPricePrecision,
}

func CreateMarketCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createmarket ",
		Short: "generate tx to create market",
		Long: "generate a tx and sign it to create market in dex blockchain." +
			"Example : createmarket [creator] [stock] [money] [priceprecision]",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			creator := cliCtx.GetFromAddress()
			msg, err := parseCreateMarketFlags(creator)
			if err != nil {
				return errors.Errorf("tx flag is error, pls see help : " +
					"$ cetcli tx market createmarket -h")
			}

			accout, err := cliCtx.GetAccount(msg.Creator)
			if err != nil {
				return errors.Errorf("Not find account with the addr : %s", msg.Creator)
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

			if !IsExistStockAndMoneySymbol(msg.Stock, msg.Money, tokens) {
				return errors.New("stock or monry is not exist in blockchain")
			}

			if msg.PricePrecision < market.MinimumTokenPricePrecision ||
				msg.PricePrecision > market.MaxTokenPricePrecision {
				return errors.Errorf("price precision out of range [%d, %d]",
					market.MinimumTokenPricePrecision, market.MaxTokenPricePrecision)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagStock, "", "The exist token symbol as stock")
	cmd.Flags().String(FlagMoney, "", "The exist token symbol as money")
	cmd.Flags().Int(FlagPricePrecision, -1, "The trading market price precision")

	for _, flag := range createMarketFlags {
		cmd.MarkFlagRequired(flag)
	}
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

func parseCreateMarketFlags(creator sdk.AccAddress) (*market.MsgCreateMarketInfo, error) {
	for _, flag := range createMarketFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, pls see help : "+
				"$ cetcli tx market createmarket", flag)
		}
	}

	msg := &market.MsgCreateMarketInfo{
		Stock:          viper.GetString(FlagStock),
		Money:          viper.GetString(FlagMoney),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Creator:        creator,
	}
	return msg, nil
}
