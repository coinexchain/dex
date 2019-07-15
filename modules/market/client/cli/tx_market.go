package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	//"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"
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

func CreateMarketCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-trading-pair ",
		Short: "generate tx to create trading pair",
		Long: `generate a tx and sign it to create trading pair in dex blockchain. 

Example : 
	cetcli tx market create-trading-pair 
	--from bob --chain-id=coinexdex 
	--stock=eth --money=cet 
	--price-precision=8 --gas 20000 --fees=1000cet`,
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			creator := cliCtx.GetFromAddress()
			msg, err := parseCreateMarketFlags(creator)
			if err != nil {
				return errors.Errorf("tx flag is error, please see help : " +
					"$ cetcli tx market createmarket -h")
			}

			//TODO we must cache the fee rates locally
			//accout, err := cliCtx.GetAccount(msg.Creator)
			//if err != nil {
			//	return errors.Errorf("Not find account with the addr : %s", msg.Creator)
			//}
			//if !accout.GetCoins().IsAllGTE(sdk.Coins{market.CreateMarketSpendCet}) {
			//	return errors.New("No have insufficient cet to create market in blockchain")
			//}

			if err := hasTokens(cliCtx, cdc, msg.Stock, msg.Money); err != nil {
				return err
			}

			if msg.PricePrecision < market.MinTokenPricePrecision ||
				msg.PricePrecision > market.MaxTokenPricePrecision {
				return errors.Errorf("price precision out of range [%d, %d]",
					market.MinTokenPricePrecision, market.MaxTokenPricePrecision)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagStock, "", "The exist token symbol as stock")
	cmd.Flags().String(FlagMoney, "", "The exist token symbol as money")
	cmd.Flags().Int(FlagPricePrecision, -1, "The trading-pair price precision")

	for _, flag := range createMarketFlags {
		cmd.MarkFlagRequired(flag)
	}
	return cmd
}

func hasTokens(cliCtx context.CLIContext, cdc *codec.Codec, tokens ...string) error {
	route := fmt.Sprintf("custom/%s/%s", asset.QuerierRoute, asset.QueryToken)
	for _, token := range tokens {
		bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(token))
		if err != nil {
			return err
		}
		fmt.Printf("token :%s\n ", token)
		if _, _, err := cliCtx.QueryWithData(route, bz); err != nil {
			fmt.Printf("route : %s\n", route)
			return err
		}
	}

	return nil
}

func parseCreateMarketFlags(creator sdk.AccAddress) (*market.MsgCreateTradingPair, error) {
	for _, flag := range createMarketFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, please see help : "+
				"$ cetcli tx market createmarket", flag)
		}
	}

	msg := &market.MsgCreateTradingPair{
		Stock:          viper.GetString(FlagStock),
		Money:          viper.GetString(FlagMoney),
		PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
		Creator:        creator,
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
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			creator := cliCtx.GetFromAddress()
			msg := market.MsgCancelTradingPair{
				Sender:        creator,
				EffectiveTime: viper.GetInt64(FlagTime),
				TradingPair:   viper.GetString(FlagSymbol),
			}

			if err := CheckCancelMarketMsg(cdc, cliCtx, msg); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagSymbol, "btc/cet", "The market trading-pair")
	cmd.Flags().Int64(FlagTime, -1, "The block height")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagTime)

	return cmd
}

func CheckCancelMarketMsg(cdc *codec.Codec, cliCtx context.CLIContext, msg market.MsgCancelTradingPair) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(market.NewQueryMarketParam(msg.TradingPair))
	if err != nil {
		return err
	}
	query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryMarket)
	res, _, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		return err
	}

	var msgMarket market.QueryMarketInfo
	if err := cdc.UnmarshalJSON(res, &msgMarket); err != nil {
		return err
	}

	if !bytes.Equal(msgMarket.Creator, msg.Sender) {
		return errors.Errorf("Not match sender, actual : %s, expect : %s\n", msg.Sender, msgMarket.Creator)
	}

	return nil
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
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc) //.WithAccountDecoder(cdc)

			creator := cliCtx.GetFromAddress()
			msg := market.MsgModifyPricePrecision{
				Sender:         creator,
				TradingPair:    viper.GetString(FlagSymbol),
				PricePrecision: byte(viper.GetInt(FlagPricePrecision)),
			}

			if err := CheckModifyPricePrecision(msg); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagSymbol, "btc/cet", "The market trading-pair")
	cmd.Flags().Int(FlagPricePrecision, 8, "The trading-pair price precision")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagPricePrecision)
	return cmd
}

func CheckModifyPricePrecision(msg market.MsgModifyPricePrecision) error {
	if len(strings.Split(msg.TradingPair, market.SymbolSeparator)) != 2 {
		return errors.Errorf("the invalid trading pair : %s ", viper.GetString(FlagSymbol))
	}
	if msg.PricePrecision < 0 || msg.PricePrecision > sdk.Precision {
		return errors.Errorf("invalid price precision : %d, expect [0, 18]", msg.PricePrecision)
	}
	return nil
}
