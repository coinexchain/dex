package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/asset"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	registerQueryRoutes(cliCtx, r, cdc, storeName)
	registerTXRoutes(cliCtx, r, cdc)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc("/asset/tokens/{symbol}", QueryTokenRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens", QueryTokensRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/whitelist", QueryWhitelistRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/addresses", QueryForbiddenAddrRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/asset/tokens/reserved/symbols", QueryReservedSymbolsRequestHandlerFn(storeName, cdc, cliCtx)).Methods("GET")
}

// QueryTokenRequestHandlerFn - query assetREST Handler
func QueryTokenRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryToken)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, asset.BaseToken{})
			return
		}

		var token asset.Token
		cdc.MustUnmarshalJSON(res, &token)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

// QueryTokensRequestHandlerFn - query assetREST Handler
func QueryTokensRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryTokenList)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			res = []byte("[]")
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// QueryWhitelistRequestHandlerFn - query assetREST Handler
func QueryWhitelistRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		bz, err := cdc.MarshalJSON(asset.NewQueryWhitelistParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryWhitelist)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			res = []byte("[]")
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// QueryForbiddenAddrRequestHandlerFn - query assetREST Handler
func QueryForbiddenAddrRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		bz, err := cdc.MarshalJSON(asset.NewQueryForbiddenAddrParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryForbiddenAddr)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			res = []byte("[]")
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// QueryReservedSymbolsRequestHandlerFn - query assetREST Handler
func QueryReservedSymbolsRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryReservedSymbols)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			res = []byte("[]")
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
