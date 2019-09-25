package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/authx/client/cli"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keeper.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}

// query accountREST Handler
func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		all := struct {
			C sdk.Coins         `json:"coins"`
			L authx.LockedCoins `json:"locked_coins"`
		}{sdk.Coins{}, authx.LockedCoins{}}
		accRetriever := auth.NewAccountRetriever(cliCtx)
		acc, height, err := accRetriever.GetAccountWithHeight(addr)
		if err != nil {
			if err := accRetriever.EnsureExists(addr); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				rest.PostProcessResponse(w, cliCtx, all)
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		all.C = acc.GetCoins()

		cliCtx = cliCtx.WithHeight(height)
		aux, err := cli.GetAccountX(cliCtx, addr)
		if err == nil {
			all.L = aux.GetAllLockedCoins()
		}
		rest.PostProcessResponse(w, cliCtx, all)
	}
}
