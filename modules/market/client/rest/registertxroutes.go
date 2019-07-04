package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/create-gte-order", createGTEOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/trading-pair", createMarketHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/create-ioc-order", createIOCOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/cancel-order", cancelOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/cancel-trading-pair", cancelMarketHandlerFn(cdc, cliCtx)).Methods("POST")
	registerQueryRoutes(cliCtx, r, cdc)
}
