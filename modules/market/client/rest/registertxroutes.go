package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/create-gte-order", createGTEOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/create-market", createMarketHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/create-ioc-order", createIOCOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/market-info/{symbol}", queryMarketHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/order-info/{order-id}", queryOrderInfoHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/user-order-list/{address}", queryUserOrderListHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/cancel-order", cancelOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/cancel-market", cancelMarketHandlerFn(cdc, cliCtx)).Methods("POST")
}
