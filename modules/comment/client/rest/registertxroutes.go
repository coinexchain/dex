package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerTXRoutes(cliCtx, r, cdc)
}

func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/comment/new-thread", createNewThreadHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/comment/followup-comment", createFollowupCommentHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/comment/reward-comments", createRewardCommentsHandlerFn(cdc, cliCtx)).Methods("POST")
}
