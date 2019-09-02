package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func queryMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stock := vars["stock"]
		money := vars["money"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := queryMarketInfo(cdc, cliCtx, stock+types.SymbolSeparator+money)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var queryInfo keepers.QueryMarketInfo
		if err := cdc.UnmarshalJSON(res, &queryInfo); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, queryInfo)
	}
}

func queryMarketsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := queryMarketInfos(cdc, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var queryInfo keepers.MarketInfoList
		if err := cdc.UnmarshalJSON(res, &queryInfo); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, queryInfo)
	}
}

func queryMarketInfos(cdc *codec.Codec, cliCtx context.CLIContext) ([]byte, int64, error) {
	query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarkets)
	res, height, err := cliCtx.QueryWithData(query, nil)
	if err != nil {
		return nil, height, err
	}
	return res, height, nil
}

func queryMarketInfo(cdc *codec.Codec, cliCtx context.CLIContext, symbol string) ([]byte, int64, error) {
	bz, err := cdc.MarshalJSON(keepers.NewQueryMarketParam(symbol))
	if err != nil {
		return nil, cliCtx.Height, err
	}

	query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryMarket)
	res, height, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		return nil, height, err
	}
	return res, height, nil
}

func queryOrderInfoHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order-id"]

		if len(strings.Split(orderID, types.OrderIDSeparator)) != types.OrderIDPartsNum {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		addr := strings.Split(orderID, types.OrderIDSeparator)[0]
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		bz, err := cdc.MarshalJSON(keepers.NewQueryOrderParam(orderID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryOrder)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryUserOrderListHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]

		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order id")
			return
		}

		bz, err := cdc.MarshalJSON(keepers.QueryUserOrderList{User: addr})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryUserOrders)
		fmt.Println(route)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
