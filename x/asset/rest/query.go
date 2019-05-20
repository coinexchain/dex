package rest

import (
	"fmt"
	"net/http"

	"github.com/coinexchain/dex/x/asset"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(
		"/asset/token/{symbol}",
		QueryTokenRequestHandlerFn(storeName, cdc, cliCtx),
	).Methods("GET")

}

// query assetREST Handler
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
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var token asset.Token
		cdc.MustUnmarshalJSON(res, &token)

		rest.PostProcessResponse(w, cdc, token, cliCtx.Indent)
	}
}
