package rest

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/bank/accounts/{address}/transfers", sendTxRequestHandlerFn(cliCtx.Codec, cliCtx)).Methods("POST")
	r.HandleFunc("/bank/accounts/memo", sendRequestHandlerFn(cliCtx.Codec, cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/bank/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}
