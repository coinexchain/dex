package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/coinexchain/dex/client/restutil"
)

func queryMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
		param := keepers.NewQueryMarketParam(vars["stock"] + types.SymbolSeparator + vars["money"])
		restutil.RestQuery(cdc, cliCtx, w, r, query, param, nil)
	}
}

func queryMarketsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarkets)
		restutil.RestQuery(cdc, cliCtx, w, r, query, nil, nil)
	}
}

func queryOrderInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := checkOrderID(vars["order-id"]); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}
		param := keepers.NewQueryOrderParam(vars["order-id"])
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
		restutil.RestQuery(cdc, cliCtx, w, r, route, param, nil)
	}
}

func checkOrderID(orderID string) error {
	if len(strings.Split(orderID, types.OrderIDSeparator)) != types.OrderIDPartsNum {
		return errors.Errorf("Invalid order id")
	}
	addr := strings.Split(orderID, types.OrderIDSeparator)[0]
	if _, err := sdk.AccAddressFromBech32(addr); err != nil {
		return errors.Errorf("Invalid order id")
	}
	return nil
}

func queryUserOrderListHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if _, err := sdk.AccAddressFromBech32(vars["address"]); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}
		param := keepers.QueryUserOrderList{User: vars["address"]}
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
		restutil.RestQuery(cdc, cliCtx, w, r, route, param, nil)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
		restutil.RestQuery(nil, cliCtx, w, r, route, nil, nil)
	}
}
