package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/staking/pool", poolHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/staking/parameters", paramsHandlerFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query the pool information
func poolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restutil.RestQuery(nil, cliCtx, w, r, "custom/stakingx/pool", nil, nil)
	}
}

// HTTP request handler to query the staking params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, staking.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}
