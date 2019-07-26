package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerTXRoutes(cliCtx, r, cdc)
	registerQueryRoutes(cliCtx, r, cdc)
}

func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/bancorlite/bancor-init", bancorInitHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/bancorlite/bancor-trade", bancorTradeHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/bancorlite/bancor-cancel", bancorCancelHandlerFn(cdc, cliCtx)).Methods("POST")
}
