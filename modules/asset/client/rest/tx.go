package rest

import (
	"fmt"
	"net/http"

	"github.com/coinexchain/dex/client/restutil"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/asset/internal/types"
)

const (
	symbol = "symbol"
)

// registerTXRoutes -
func registerTXRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/asset/tokens", issueRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/ownerships", transferOwnerRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/mints", mintTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/burns", burnTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/forbids", forbidTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/unforbids", unForbidTokenHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/whitelist", addWhitelistHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/unforbidden/whitelist", removeWhitelistHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/forbidden/addresses", forbidAddrHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/unforbidden/addresses", unForbidAddrHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/asset/tokens/{symbol}/infos", modifyTokenInfoHandlerFn(cdc, cliCtx)).Methods("POST")
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

		amt, ok := sdk.NewIntFromString(req.TotalSupply)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidTokenSupply(req.TotalSupply).Error())
			return
		}

		bz, err := cdc.MarshalJSON(types.NewQueryAssetParams(req.Symbol))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryToken)
		if res, _, _ := cliCtx.QueryWithData(route, bz); res != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrDuplicateTokenSymbol(symbol).Error())
			return
		}

		msg := types.NewMsgIssueToken(req.Name, req.Symbol, amt, owner,
			req.Mintable, req.Burnable, req.AddrForbiddable, req.TokenForbiddable, req.URL, req.Description, req.Identity)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// transferOwnershipRequestHandlerFn - http request handler to transfer token owner ship.
func transferOwnerRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(transferOwnerReq))
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

		symbol := getSymbol(r)
		amt, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidTokenMintAmt(req.Amount).Error())
			return
		}
		msg := types.NewMsgMintToken(symbol, amt, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

		symbol := getSymbol(r)
		amt, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidTokenBurnAmt(req.Amount).Error())
			return
		}
		msg := types.NewMsgBurnToken(symbol, amt, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

		symbol := getSymbol(r)

		msg := types.NewMsgForbidToken(symbol, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

		symbol := getSymbol(r)

		msg := types.NewMsgUnForbidToken(symbol, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// addWhitelistHandlerFn - http request handler to add whitelist.
func addWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addrListReq
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

		symbol := getSymbol(r)

		msg := types.NewMsgAddTokenWhitelist(symbol, owner, req.AddrList)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// removeWhitelistHandlerFn - http request handler to add whitelist.
func removeWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addrListReq
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

		symbol := getSymbol(r)

		msg := types.NewMsgRemoveTokenWhitelist(symbol, owner, req.AddrList)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// forbidAddrHandlerFn - http request handler to forbid addresses.
func forbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addrListReq
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

		symbol := getSymbol(r)

		msg := types.NewMsgForbidAddr(symbol, owner, req.AddrList)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// unForbidAddrHandlerFn - http request handler to unforbid addresses.
func unForbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addrListReq
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

		symbol := getSymbol(r)

		msg := types.NewMsgUnForbidAddr(symbol, owner, req.AddrList)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// modifyTokenInfoHandlerFn - http request handler to modify token url.
func modifyTokenInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req modifyTokenInfoReq
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

		url := types.DoNotModifyTokenInfo
		if req.URL != nil {
			url = *req.URL
		}

		description := types.DoNotModifyTokenInfo
		if req.Description != nil {
			description = *req.Description
		}

		identity := types.DoNotModifyTokenInfo
		if req.Identity != nil {
			identity = *req.Identity
		}

		symbol := getSymbol(r)
		msg := types.NewMsgModifyTokenInfo(symbol, url, description, identity, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func getSymbol(r *http.Request) string {
	vars := mux.Vars(r)
	return vars[symbol]
}
