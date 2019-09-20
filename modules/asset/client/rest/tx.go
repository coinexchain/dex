package rest

import (
	"fmt"
	"net/http"

	"github.com/coinexchain/dex/client/restutil"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

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
	checker := func(cdc *codec.Codec, cliCtx context.CLIContext, req restutil.RestReq) error {
		symbol := req.(*issueReq).Symbol
		bz, err := cdc.MarshalJSON(types.NewQueryAssetParams(symbol))
		if err != nil {
			return err
		}
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryToken)
		if res, _, _ := cliCtx.QueryWithData(route, bz); res != nil {
			return types.ErrDuplicateTokenSymbol(symbol)
		}
		return nil
	}
	return restutil.NewRestHandlerBuilder(cdc, cliCtx, new(issueReq)).Build(checker)
}

// transferOwnershipRequestHandlerFn - http request handler to transfer token owner ship.
func transferOwnerRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(transferOwnerReq))
}

// mintTokenHandlerFn - http request handler to mint token.
func mintTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(mintTokenReq))
}

// burnTokenHandlerFn - http request handler to burn token.
func burnTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(burnTokenReq))
}

// forbidTokenHandlerFn - http request handler to forbid token.
func forbidTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(forbidTokenReq))
}

// unForbidTokenHandlerFn - http request handler to unforbid token.
func unForbidTokenHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(unForbidTokenReq))
}

// addWhitelistHandlerFn - http request handler to add whitelist.
func addWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(addWhiteListReq))
}

// removeWhitelistHandlerFn - http request handler to add whitelist.
func removeWhitelistHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(removeWhiteListReq))
}

// forbidAddrHandlerFn - http request handler to forbid addresses.
func forbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(forbidAddrReq))
}

// unForbidAddrHandlerFn - http request handler to unforbid addresses.
func unForbidAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(unforbidAddrReq))
}

// modifyTokenInfoHandlerFn - http request handler to modify token url.
func modifyTokenInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(modifyTokenInfoReq))
}
