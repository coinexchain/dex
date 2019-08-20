package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/asset/internal/types"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	registerQueryRoutes(cliCtx, r, cdc, storeName)
	registerTXRoutes(cliCtx, r, cdc)
}

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
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		bz, err := cdc.MarshalJSON(types.NewQueryAssetParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryToken)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, types.BaseToken{})
			return
		}

		var token types.Token
		cdc.MustUnmarshalJSON(res, &token)

		rest.PostProcessResponse(w, cliCtx, token)
	}
}

// QueryTokensRequestHandlerFn - query assetREST Handler
func QueryTokensRequestHandlerFn(
	storeName string, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryTokenList)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
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

		bz, err := cdc.MarshalJSON(types.NewQueryWhitelistParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryWhitelist)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
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

		bz, err := cdc.MarshalJSON(types.NewQueryForbiddenAddrParams(symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryForbiddenAddr)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			res = []byte("[]")
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// QueryReservedSymbolsRequestHandlerFn - query assetREST Handler
func QueryReservedSymbolsRequestHandlerFn(
	storeName string, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryReservedSymbols)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			res = []byte("[]")
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the asset params values
func QueryParamsHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
