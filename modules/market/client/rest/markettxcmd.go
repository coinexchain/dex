package rest

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/coinexchain/dex/modules/market/client/cli"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/market"
)

// SendReq defines the properties of a send request's body.
type createMarketReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Stock          string       `json:"stock"`
	Money          string       `json:"money"`
	PricePrecision int          `json:"price_precision"`
}

type cancelMarketReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Symbol  string       `json:"symbol"`
	Time    int64        `json:"time"`
}

func createMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createMarketReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence, err := cliCtx.GetAccountSequence(creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Don't get sequence from blockchain.")
			return
		}
		req.BaseReq.Sequence = sequence

		msg := market.NewMsgCreateMarketInfo(req.Stock, req.Money, creator, byte(req.PricePrecision))
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func queryMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		if len(strings.Split(symbol, market.SymbolSeparator)) != 2 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "The invalid symbol")
			return
		}

		res, err := queryMarketInfo(cdc, cliCtx, symbol)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryMarketInfo(cdc *codec.Codec, cliCtx context.CLIContext, symbol string) ([]byte, error) {
	bz, err := cdc.MarshalJSON(market.NewQueryMarketParam(symbol))
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryMarket)
	res, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func cancelMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var req cancelMarketReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, _ := sdk.AccAddressFromBech32(req.BaseReq.From)
		msg := market.MsgCancelMarket{
			Sender:        sender,
			Symbol:        req.Symbol,
			EffectiveTime: req.Time,
		}

		if err := cli.CheckCancelMarketMsg(cdc, cliCtx, msg); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
