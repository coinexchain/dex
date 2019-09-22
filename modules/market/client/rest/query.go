package rest

import (
	"fmt"
	"net/http"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/coinexchain/dex/client/restutil"
)

func queryMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
		param := keepers.NewQueryMarketParam(vars["stock"] + types.SymbolSeparator + vars["money"])
		restutil.RestQuery(cdc, cliCtx, w, r, query, param, nil)
	}
}

func queryMarketsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarkets)
		restutil.RestQuery(cdc, cliCtx, w, r, query, nil, nil)
	}
}

func queryOrderInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := keepers.NewQueryOrderParam(vars["order-id"])
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
		restutil.RestQuery(cdc, cliCtx, w, r, route, param, nil)
	}
}

func queryUserOrderListHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := keepers.QueryUserOrderList{User: vars["address"]}
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
		restutil.RestQuery(cdc, cliCtx, w, r, route, param, nil)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}
