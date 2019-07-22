package rest

import (
	"net/http"

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

type (
	// issueReq defines the properties of a issue token request's body
	issueReq struct {
		BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
		Name             string       `json:"name" yaml:"name"`
		Symbol           string       `json:"symbol" yaml:"symbol"`
		TotalSupply      int64        `json:"total_supply" yaml:"total_supply"`
		Mintable         bool         `json:"mintable" yaml:"mintable"`
		Burnable         bool         `json:"burnable" yaml:"burnable"`
		AddrForbiddable  bool         `json:"addr_forbiddable" yaml:"addr_forbiddable"`
		TokenForbiddable bool         `json:"token_forbiddable" yaml:"token_forbiddable"`
		URL              string       `json:"url" yaml:"url"`
		Description      string       `json:"description" yaml:"description"`
	}

	// transferOwnerReq defines the properties of a transfer ownership request's body.
	transferOwnerReq struct {
		BaseReq  rest.BaseReq   `json:"base_req" yaml:"base_req"`
		NewOwner sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
	}

	// mintTokenReq defines the properties of a mint token request's body.
	mintTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  int64        `json:"amount" yaml:"amount"`
	}

	// burnTokenReq defines the properties of a burn token request's body.
	burnTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount  int64        `json:"amount" yaml:"amount"`
	}

	// forbidTokenReq defines the properties of a forbid token request's body.
	forbidTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
	// unforbidTokenReq defines the properties of a unforbid token request's body.
	unForbidTokenReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
	// addrListReq defines the properties of a whitelist or forbidden addr request's body.
	addrListReq struct {
		BaseReq  rest.BaseReq     `json:"base_req" yaml:"base_req"`
		AddrList []sdk.AccAddress `json:"addr_list" yaml:"addr_list"`
	}
	// modifyTokenInfoReq defines the properties of a modify token info request's body.
	modifyTokenInfoReq struct {
		BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
		URL         string       `json:"url" yaml:"url"`
		Description string       `json:"description" yaml:"description"`
	}
)

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

		msg := types.NewMsgIssueToken(req.Name, req.Symbol, req.TotalSupply, owner,
			req.Mintable, req.Burnable, req.AddrForbiddable, req.TokenForbiddable, req.URL, req.Description)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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
		msg := types.NewMsgTransferOwnership(symbol, original, req.NewOwner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
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

		msg := types.NewMsgMintToken(symbol, req.Amount, owner)
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

		msg := types.NewMsgBurnToken(symbol, req.Amount, owner)
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

		symbol := getSymbol(r)

		msg := types.NewMsgModifyTokenInfo(symbol, req.URL, req.Description, owner)
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
