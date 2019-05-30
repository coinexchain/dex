package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market"
	"github.com/spf13/viper"
	"strings"

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
			"Example : " +
			" cetcli tx market createmarket " +
			"--from bob --chain-id=coinexdex " +
			"--stock=eth --money=cet " +
			"--price-precision=8 --gas 20000 ",
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

			if err := hasTokens(cliCtx, cdc, queryRoute, msg.Stock, msg.Money); err != nil {
				return err
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

func QueryMarketCmd(cdc *codec.Codec) *cobra.Command {

	return &cobra.Command{
		Use:   "marketinfo",
		Short: "query market info",
		Long: "cetcli query market marketinfo [symbol]" +
			"Example : " +
			"cetcli query market " +
			"marketinfo eth/cet " +
			"--trust-node=true",
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
			query := fmt.Sprintf("custom/%s/%s", market.MarketKey, market.QueryMarket)
			res, err := cliCtx.QueryWithData(query, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}
