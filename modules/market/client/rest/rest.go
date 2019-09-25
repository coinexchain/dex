package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTXRoutes(cliCtx, r, cdc)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/trading-pairs/{stock}/{money}", queryMarketHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/exist-trading-pairs", queryMarketsHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/orders/{order-id}", queryOrderInfoHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/orders/account/{address}", queryUserOrderListHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/market/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/gte-orders", createGTEOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/trading-pairs", createMarketHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/ioc-orders", createIOCOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/cancel-order", cancelOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/cancel-trading-pair", cancelMarketHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/price-precision", modifyTradingPairPricePrecision(cdc, cliCtx)).Methods("POST")
}
