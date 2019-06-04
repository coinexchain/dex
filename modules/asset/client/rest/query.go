package rest

import (
	"fmt"
	"net/http"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	registerQueryRoutes(cliCtx, r, cdc, storeName)
	registerTXRoutes(cliCtx, r, cdc)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(
		"/asset/tokens/{symbol}",
		QueryTokenRequestHandlerFn(storeName, cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/asset/tokens",
		QueryTokensRequestHandlerFn(storeName, cdc, cliCtx),
	).Methods("GET")

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
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryToken)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if len(res) == 0 {
			rest.PostProcessResponse(w, cdc, asset.BaseToken{}, cliCtx.Indent)
			return
		}

		var token asset.Token
		cdc.MustUnmarshalJSON(res, &token)

		rest.PostProcessResponse(w, cdc, token, cliCtx.Indent)
	}
}

// QueryTokensRequestHandlerFn - query assetREST Handler
func QueryTokensRequestHandlerFn(
	storeName string, cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", storeName, asset.QueryTokenList)
		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			res = []byte("[]")
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
