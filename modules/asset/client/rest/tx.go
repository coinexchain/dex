package rest

import (
	"github.com/coinexchain/dex/modules/asset"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// registerTXRoutes -
func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/asset/tokens", issueRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}

// SendReq defines the properties of a send request's body.
type issueReq struct {
	BaseReq         rest.BaseReq `json:"basereq"`
	Name            string       `json:"name"`
	Symbol          string       `json:"symbol"`
	TotalSupply     int64        `json:"total_supply"`
	Mintable        bool         `json:"mintable"`
	Burnable        bool         `json:"burnable"`
	AddrFreezeable  bool         `json:"addrfreezeable"`
	TokenFreezeable bool         `json:"tokenfreezeable"`
}

// issueRequestHandlerFn - http request handler to issue new token.
func issueRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := asset.NewMsgIssueToken(req.Name, req.Symbol, req.TotalSupply, owner,
			req.Mintable, req.Burnable, req.AddrFreezeable, req.TokenFreezeable)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
