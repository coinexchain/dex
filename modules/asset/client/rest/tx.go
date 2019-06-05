package rest

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	symbol = "symbol"
)

// registerTXRoutes -
func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/asset/tokens", issueRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/ownerships", symbol), transferOwnerRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/mints", symbol), mintTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/burns", symbol), burnTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/forbids", symbol), forbidTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/unforbids", symbol), unForbidTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/forbidden/whitelists", symbol), addWhitelistHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/unforbidden/whitelists", symbol), removeWhitelistHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/forbidden/addresses", symbol), forbidAddrHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}/unforbidden/addresses", symbol), unForbidAddrHandlerFn(cdc, cliCtx)).Methods("POST")
}

// issueReq defines the properties of a issue token request's body.
type issueReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	Name             string       `json:"name"`
	Symbol           string       `json:"symbol"`
	TotalSupply      int64        `json:"total_supply"`
	Mintable         bool         `json:"mintable"`
	Burnable         bool         `json:"burnable"`
	AddrForbiddable  bool         `json:"addr_forbiddable"`
	TokenForbiddable bool         `json:"token_forbiddable"`
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
			req.Mintable, req.Burnable, req.AddrForbiddable, req.TokenForbiddable)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// transferOwnerReq defines the properties of a transfer ownership request's body.
type transferOwnerReq struct {
	BaseReq  rest.BaseReq   `json:"base_req"`
	NewOwner sdk.AccAddress `json:"new_owner"`
}

// transferOwnershipRequestHandlerFn - http request handler to transfer token owner ship.
func transferOwnerRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferOwnerReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		original, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		symbol := vars["symbol"]
		msg := asset.NewMsgTransferOwnership(symbol, original, req.NewOwner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// mintTokenReq defines the properties of a mint token request's body.
type mintTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount  int64        `json:"amount"`
}

// mintTokenHandlerFn - http request handler to mint token.
func mintTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintTokenReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgMintToken(symbol, req.Amount, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// burnTokenReq defines the properties of a burn token request's body.
type burnTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount  int64        `json:"amount"`
}

// burnTokenHandlerFn - http request handler to burn token.
func burnTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnTokenReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgBurnToken(symbol, req.Amount, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// forbidTokenReq defines the properties of a forbid token request's body.
type forbidTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

// forbidTokenHandlerFn - http request handler to forbid token.
func forbidTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req forbidTokenReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgForbidToken(symbol, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// unforbidTokenReq defines the properties of a unforbid token request's body.
type unForbidTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

// unForbidTokenHandlerFn - http request handler to unforbid token.
func unForbidTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req unForbidTokenReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgUnForbidToken(symbol, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// addWhitelistReq defines the properties of a add whitelist request's body.
type addWhitelistReq struct {
	BaseReq   rest.BaseReq     `json:"base_req"`
	Whitelist []sdk.AccAddress `json:"whitelist"`
}

// addWhitelistHandlerFn - http request handler to add whitelist.
func addWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addWhitelistReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgAddTokenWhitelist(symbol, owner, req.Whitelist)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// removeWhitelistReq defines the properties of a remove whitelist request's body.
type removeWhitelistReq struct {
	BaseReq   rest.BaseReq     `json:"base_req"`
	Whitelist []sdk.AccAddress `json:"whitelist"`
}

// removeWhitelistHandlerFn - http request handler to add whitelist.
func removeWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req removeWhitelistReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgRemoveTokenWhitelist(symbol, owner, req.Whitelist)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// forbidAddrReq defines the properties of a forbid addr request's body.
type forbidAddrReq struct {
	BaseReq rest.BaseReq     `json:"base_req"`
	Addr    []sdk.AccAddress `json:"addr"`
}

// forbidAddrHandlerFn - http request handler to forbid addresses.
func forbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req forbidAddrReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgForbidAddr(symbol, owner, req.Addr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// unForbidAddrReq defines the properties of a unforbid addr request's body.
type unForbidAddrReq struct {
	BaseReq rest.BaseReq     `json:"base_req"`
	Addr    []sdk.AccAddress `json:"addr"`
}

// unForbidAddrHandlerFn - http request handler to unforbid addresses.
func unForbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req unForbidAddrReq
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

		vars := mux.Vars(r)
		symbol := vars[symbol]

		msg := asset.NewMsgUnForbidAddr(symbol, owner, req.Addr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
