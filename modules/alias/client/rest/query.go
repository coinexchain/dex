package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	"github.com/coinexchain/dex/modules/authx/client/restutil"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/alias/address-of-alias/{alias}", queryAddressHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/alias/aliases-of-address/{address}", queryAliasesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/alias/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func queryAddressHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
		param := &keepers.QueryAliasInfoParam{Alias: vars["alias"], QueryOp: keepers.GetAddressFromAlias}

		restutil.RestQuery(cdc, cliCtx, w, r, query, param)
	}
}

func queryAliasesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryAliasInfo)
		acc, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		param := &keepers.QueryAliasInfoParam{Owner: acc, QueryOp: keepers.ListAliasOfAccount}
		restutil.RestQuery(cdc, cliCtx, w, r, query, param)
	}
}

// HTTP request handler to query the alias params values
func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil)
	}
}
