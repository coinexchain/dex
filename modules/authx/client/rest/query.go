package rest

import (
	"net/http"

	"github.com/coinexchain/dex/modules/authx/client/cli"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/coinexchain/dex/modules/authx"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/auth/accounts/{address}",
		QueryAccountRequestHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/bank/balances/{address}",
		QueryBalancesRequestHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		accRetriever := authtypes.NewAccountRetriever(cliCtx)
		if err = accRetriever.EnsureExists(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		acc, err := accRetriever.GetAccount(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		aux, err := cli.GetAccountX(cliCtx, addr)
		if err != nil {
			aux = authx.AccountX{}
		}

		mix := authx.NewAccountMix(acc, aux)

		rest.PostProcessResponse(w, cliCtx, mix)
	}
}

// query accountREST Handler
func QueryBalancesRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		accRetriever := authtypes.NewAccountRetriever(cliCtx)
		if err = accRetriever.EnsureExists(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		acc, err := accRetriever.GetAccount(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		lockedCoins := make(authx.LockedCoins, 0)
		aux, err := cli.GetAccountX(cliCtx, addr)
		if err == nil {
			lockedCoins = aux.GetAllLockedCoins()
		}

		all := struct {
			C sdk.Coins         `json:"coins"`
			L authx.LockedCoins `json:"locked_coins"`
		}{acc.GetCoins(), lockedCoins}

		rest.PostProcessResponse(w, cliCtx, all)
	}
}
