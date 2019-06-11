package rest

import (
	clientx "github.com/coinexchain/dex/client"
	"github.com/coinexchain/dex/modules/authx"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"net/http"

	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(
		"/auth/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/bank/balances/{address}",
		QueryBalancesRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = cliCtx.EnsureAccountExistsFromAddr(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		acc, err := cliCtx.GetAccount(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		aux, err := clientx.GetAccountX(cliCtx, addr)
		if err != nil {
			rest.PostProcessResponse(w, cdc, acc, cliCtx.Indent)
			return
		}

		mix := authx.NewAccountMix(acc, aux)

		rest.PostProcessResponse(w, cdc, mix, cliCtx.Indent)
	}
}

// query accountREST Handler
func QueryBalancesRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = cliCtx.EnsureAccountExistsFromAddr(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		acc, err := cliCtx.GetAccount(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		lockedCoins := make(authx.LockedCoins, 0)
		aux, err := clientx.GetAccountX(cliCtx, addr)
		if err == nil {
			lockedCoins = aux.GetAllLockedCoins()
		}

		all := struct {
			C sdk.Coins         `json:"coins"`
			L authx.LockedCoins `json:"locked_coins"`
		}{acc.GetCoins(), lockedCoins}

		rest.PostProcessResponse(w, cdc, all, cliCtx.Indent)
	}
}
