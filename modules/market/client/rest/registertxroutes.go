package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/market/creategteorder", createGTEOrderHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/market/createmarket", createMarketHandlerFn(cdc, cliCtx)).Methods("POST")
}
