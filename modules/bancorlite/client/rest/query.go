package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/bancorlite/pools/{symbol}", queryBancorInfoHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/bancorlite/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}

// format: barcorlite/pools/btc-cet
func queryBancorInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryBancorInfo)
		symbol := strings.Replace(vars["symbol"], "-", "/", 1)
		param := &keepers.QueryBancorInfoParam{Symbol: symbol}
		restutil.RestQuery(cdc, cliCtx, w, r, query, param, nil)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}
