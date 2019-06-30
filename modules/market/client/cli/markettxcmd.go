package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

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

func CreateMarketCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trading-pair ",
		Short: "generate tx to create trading pair",
		Long: "generate a tx and sign it to create trading pair in dex blockchain. \n" +
			"Example : " +
			" cetcli tx market trading-pair " +
			"--from bob --chain-id=coinexdex " +
			"--stock=eth --money=cet " +
			"--price-precision=8 --gas 20000 --fees=1000cet ",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

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

			if err := hasTokens(cliCtx, cdc, queryRoute, msg.Stock, msg.Money); err != nil {
				return err
			}

			if msg.PricePrecision < market.MinTokenPricePrecision ||
				msg.PricePrecision > market.MaxTokenPricePrecision {
				return errors.Errorf("price precision out of range [%d, %d]",
					market.MinTokenPricePrecision, market.MaxTokenPricePrecision)
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

func hasTokens(cliCtx context.CLIContext, cdc *codec.Codec, queryRoute string, tokens ...string) error {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryToken)
	for _, token := range tokens {
		bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(token))
		if err != nil {
			return err
		}
		fmt.Printf("token :%s\n ", token)
		if _, err := cliCtx.QueryWithData(route, bz); err != nil {
			fmt.Printf("route : %s\n", route)
			return err
		}
	}

	return nil
}

func parseCreateMarketFlags(creator sdk.AccAddress) (*market.MsgCreateMarketInfo, error) {
	for _, flag := range createMarketFlags {
		if viper.Get(flag) == nil {
			return nil, fmt.Errorf("--%s flag is a noop, please see help : "+
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

func QueryMarketCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "trading-pair",
		Short: "query trading-pair info in blockchain",
		Long: "query trading-pair info in blockchain. \n" +
			"Example : " +
			"cetcli query market " +
			"trading-pair eth/cet " +
			"--trust-node=true --chain-id=coinexdex",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			if len(strings.Split(args[0], market.SymbolSeparator)) != 2 {
				return errors.Errorf("symbol illegal : %s, For example : eth/cet.", args[0])
			}

			bz, err := cdc.MarshalJSON(market.NewQueryMarketParam(args[0]))
			if err != nil {
				return err
			}
			query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryMarket)
			res, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

func CancelMarket(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-trading-pair",
		Short: "cancel trading-pair in blockchain",
		Long: "cancel trading-pair in blockchain at least a week from now. \n " +
			"Example : " +
			"cetcli tx market cancel-trading-pair " +
			"--time=1000000 --symbol=etc/cet --from=bob --chain-id=coinexdex " +
			"--gas=1000000 --fees=1000cet",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			creator := cliCtx.GetFromAddress()
			msg := market.MsgCancelMarket{
				Sender:        creator,
				EffectiveTime: viper.GetInt64(FlagTime),
				Symbol:        viper.GetString(FlagSymbol),
			}

			if err := CheckCancelMarketMsg(cdc, cliCtx, msg); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "The market symbol")
	cmd.Flags().Int64(FlagTime, -1, "The block height")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagTime)

	return cmd
}

func CheckCancelMarketMsg(cdc *codec.Codec, cliCtx context.CLIContext, msg market.MsgCancelMarket) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(market.NewQueryMarketParam(msg.Symbol))
	if err != nil {
		return err
	}
	query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryMarket)
	res, err := cliCtx.QueryWithData(query, bz)
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

func QueryWaitCancelMarkets(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wait-cancel-trading-pair",
		Short: "Query wait cancel trading-pair info in special time",
		Long: "Query wait cancel trading-pair info in special time \n" +
			"Example:" +
			"cetcli query market " +
			"wait-cancel-trading-pair --time=10000 " +
			"--trust-node=true --chain-id=coinexdex",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			time := viper.GetInt64(FlagTime)
			if time <= 0 {
				return errors.Errorf("Invalid unix time")
			}

			bz, err := cdc.MarshalJSON(market.QueryCancelMarkets{Time: time})
			if err != nil {
				return err
			}

			query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryWaitCancelMarkets)
			res, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}

			var markets []string
			if err := cdc.UnmarshalJSON(res, &markets); err != nil {
				return err
			}
			fmt.Println(markets)

			return nil
		},
	}

	cmd.Flags().Int64(FlagTime, -1, "The query block height")
	cmd.MarkFlagRequired(FlagTime)
	return cmd
}
