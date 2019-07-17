package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/comment"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/comment/get-count", queryCommentCountHandlerFn(cdc, cliCtx)).Methods("GET")
}

func queryCommentCountHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := fmt.Sprintf("custom/%s/%s", comment.StoreKey, comment.QueryCommentCount)
		res, _, err := cliCtx.QueryWithData(query, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
