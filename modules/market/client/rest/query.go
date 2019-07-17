package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"

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

		res, err := queryMarketInfo(cdc, cliCtx, stock+types.SymbolSeparator+money)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var queryInfo keepers.QueryMarketInfo
		if err := cdc.UnmarshalJSON(res, &queryInfo); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, queryInfo)
	}
}

func queryMarketInfo(cdc *codec.Codec, cliCtx context.CLIContext, symbol string) ([]byte, error) {
	bz, err := cdc.MarshalJSON(keepers.NewQueryMarketParam(symbol))
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
	res, _, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func queryOrderInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order-id"]

		if len(strings.Split(orderID, "-")) != 3 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		addr := strings.Split(orderID, "-")[0]
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		bz, err := cdc.MarshalJSON(keepers.NewQueryOrderParam(orderID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
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

		bz, err := cdc.MarshalJSON(keepers.QueryUserOrderList{User: addr})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
		fmt.Println(route)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
