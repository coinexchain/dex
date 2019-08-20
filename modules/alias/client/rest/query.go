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
		bz, err := cdc.MarshalJSON(param)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := cliCtx.QueryWithData(query, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
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
		bz, err := cdc.MarshalJSON(param)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := cliCtx.QueryWithData(query, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the alias params values
func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/alias/parameters", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
