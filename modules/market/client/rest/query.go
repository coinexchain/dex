package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/coinexchain/dex/modules/market"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/trading-pair/{stock}/{money}", queryMarketHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/order-info/{order-id}", queryOrderInfoHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/user-order-list/{address}", queryUserOrderListHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stock := vars["stock"]
		money := vars["money"]

		res, err := queryMarketInfo(cdc, cliCtx, stock+market.SymbolSeparator+money)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var queryInfo market.QueryMarketInfo
		if err := cdc.UnmarshalJSON(res, &queryInfo); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, queryInfo, cliCtx.Indent)
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

func queryOrderInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order-id"]

		if len(strings.Split(orderID, "-")) != 2 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		addr := strings.Split(orderID, "-")[0]
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		bz, err := cdc.MarshalJSON(market.NewQueryOrderParam(orderID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryOrder)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryUserOrderListHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]

		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		bz, err := cdc.MarshalJSON(market.QueryUserOrderList{User: addr})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", market.StoreKey, market.QueryUserOrders)
		fmt.Println(route)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
