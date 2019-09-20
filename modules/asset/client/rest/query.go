package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

var (
	emptyJSONObj = []byte("{}")
	emptyJSONArr = []byte("[]")
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc("/asset/tokens/{symbol}", QueryTokenRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens", QueryTokensRequestHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/whitelist", QueryWhitelistRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/addresses", QueryForbiddenAddrRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/reserved/symbols", QueryReservedSymbolsRequestHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/parameters", QueryParamsHandlerFn(storeName, cliCtx)).Methods("GET")
}

// QueryTokenRequestHandlerFn - query assetREST Handler
func QueryTokenRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryToken)
		symbol := mux.Vars(r)["symbol"]
		params := types.NewQueryAssetParams(symbol)
		restutil.RestQuery(cdc, cliCtx, w, r, route, params, emptyJSONObj)
	}
}

// QueryTokensRequestHandlerFn - query assetREST Handler
func QueryTokensRequestHandlerFn(
	storeName string, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryTokenList)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, emptyJSONArr)
	}
}

// QueryWhitelistRequestHandlerFn - query assetREST Handler
func QueryWhitelistRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryWhitelist)
		symbol := mux.Vars(r)["symbol"]
		params := types.NewQueryWhitelistParams(symbol)
		restutil.RestQuery(cdc, cliCtx, w, r, route, params, emptyJSONArr)
	}
}

// QueryForbiddenAddrRequestHandlerFn - query assetREST Handler
func QueryForbiddenAddrRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryForbiddenAddr)
		symbol := mux.Vars(r)["symbol"]
		params := types.NewQueryForbiddenAddrParams(symbol)
		restutil.RestQuery(cdc, cliCtx, w, r, route, params, emptyJSONArr)
	}
}

// QueryReservedSymbolsRequestHandlerFn - query assetREST Handler
func QueryReservedSymbolsRequestHandlerFn(
	storeName string, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryReservedSymbols)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}

// HTTP request handler to query the asset params values
func QueryParamsHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}
