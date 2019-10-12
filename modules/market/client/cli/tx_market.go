package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

const (
	FlagStock          = "stock"
	FlagMoney          = "money"
	FlagPricePrecision = "price-precision"
	FlagOrderPrecision = "order-precision"
)

var createMarketFlags = []string{
	FlagStock,
	FlagMoney,
	FlagPricePrecision,
}

func CreateMarketCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-trading-pair",
		Short: "generate tx to create trading pair",
		Long: `generate a tx and sign it to create trading pair in dex blockchain. 

Example : 
	cetcli tx market create-trading-pair  \
	--from bob --chain-id=coinexdex  \
	--stock=eth --money=cet --order-precision=8 \
	--price-precision=8 --gas 20000 --fees=1000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := getCreateMarketMsg(cdc)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(FlagStock, "", "The exist token symbol as stock")
	cmd.Flags().String(FlagMoney, "", "The exist token symbol as money")
	cmd.Flags().Int(FlagPricePrecision, 1, "The trading-pair price precision, used to"+
		" control the price accuracy of the order when token trades")
	cmd.Flags().Int(FlagOrderPrecision, 0, "To control the granularity of token trade, "+
		"the token amount of trade must be a multiple of granularity.")
	for _, flag := range createMarketFlags {
		cmd.MarkFlagRequired(flag)
	}
	return cmd
}

func getCreateMarketMsg(cdc *codec.Codec) (*types.MsgCreateTradingPair, error) {
	msg, err := parseCreateMarketFlags()
	if err != nil {
		return nil, errors.Errorf("tx flag is error, please see help : " +
			"$ cetcli tx market createmarket -h")
	}
	if err := hasTokens(cdc, msg.Stock, msg.Money); err != nil {
		return nil, err
	}
	return msg, nil
}

func hasTokens(cdc *codec.Codec, tokens ...string) error {
	route := fmt.Sprintf("custom/%s/%s", asset.QuerierRoute, asset.QueryToken)
	for _, token := range tokens {
		if err := asset.ValidateTokenSymbol(token); err != nil {
			return err
		}
		if err := cliutil.CliQuery(cdc, route, asset.NewQueryAssetParams(token)); err != nil {
			return err
		}
	}
	return nil
}

func parseCreateMarketFlags() (*types.MsgCreateTradingPair, error) {
	for _, flag := range createMarketFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, please see help : "+
				"$ cetcli tx market createmarket", flag)
		}
	}

	msg := &types.MsgCreateTradingPair{
		Stock:          viper.GetString(FlagStock),
		Money:          viper.GetString(FlagMoney),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		OrderPrecision: byte(viper.GetInt(FlagOrderPrecision)),
	}
	return msg, nil
}

func CancelMarket(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-trading-pair",
		Short: "cancel trading-pair in blockchain",
		Long: `cancel trading-pair in blockchain at least a week from now. 

Example 
	cetcli tx market cancel-trading-pair 
	--time=1000000 --trading-pair=etc/cet --from=bob --chain-id=coinexdex 
	--gas=1000000 --fees=1000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := getCancelMarketMsg(cdc)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(FlagSymbol, "btc/cet", "The market trading-pair")
	cmd.Flags().Int64(FlagTime, 100, "The trading pair on expired unix timestamp. (timestamp - time.Now() > 7days)")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagTime)

	return cmd
}

func getCancelMarketMsg(cdc *codec.Codec) (*types.MsgCancelTradingPair, error) {
	msg := types.MsgCancelTradingPair{
		EffectiveTime: viper.GetInt64(FlagTime),
		TradingPair:   viper.GetString(FlagSymbol),
	}
	query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
	if err := cliutil.CliQuery(cdc, query, keepers.NewQueryMarketParam(msg.TradingPair)); err != nil {
		return nil, err
	}
	return &msg, nil
}

func ModifyTradingPairPricePrecision(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-price-precision",
		Short: "Modify the price precision of the trading pair ",
		Long: `Modify the price precision of the trading pair in the dex.

Example: 
	cetcli tx market modify-price-precision --trading-pair=etc/cet
	--price-precision=9 --from=bob --chain-id=coinexdex 
	--gas=10000000 --fees=10000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := getModifyTradingPairPricePrecisionMsg(cdc)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(FlagSymbol, "btc/cet", "The market trading-pair")
	cmd.Flags().Int(FlagPricePrecision, 8, "The trading-pair price precision")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagPricePrecision)
	return cmd
}

func getModifyTradingPairPricePrecisionMsg(cdc *codec.Codec) (*types.MsgModifyPricePrecision, error) {
	msg := types.MsgModifyPricePrecision{
		TradingPair:    viper.GetString(FlagSymbol),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
	}
	query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
	if err := cliutil.CliQuery(cdc, query, keepers.NewQueryMarketParam(msg.TradingPair)); err != nil {
		return nil, err
	}
	return &msg, nil
}
