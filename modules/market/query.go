package market

import (
	"fmt"

	"github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryMarket            = "market-info"
	QueryOrder             = "order-info"
	QueryUserOrders        = "user-order-list"
	QueryWaitCancelMarkets = "wait-cancel-markets"
)

// creates a querier for asset REST endpoints
func NewQuerier(mk Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req types.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryMarket:
			return queryMarket(ctx, req, mk)
		case QueryOrder:
			return queryOrder(ctx, req, mk)
		case QueryUserOrders:
			return queryUserOrderList(ctx, req, mk)
		case QueryWaitCancelMarkets:
			return queryWaitCancelMarkets(ctx, req, mk)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
		}
	}
}

type QueryMarketParam struct {
	Symbol string
}

func NewQueryMarketParam(symbol string) QueryMarketParam {
	return QueryMarketParam{
		Symbol: symbol,
	}
}

func queryMarket(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryMarketParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse param: %s", err))
	}

	info, err := mk.GetMarketInfo(ctx, param.Symbol)
	if err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeInvalidSymbol, "may be the market have deleted or not exist")
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, info)
	if err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}

type QueryOrderParam struct {
	OrderID string
}

func NewQueryOrderParam(orderID string) QueryOrderParam {
	return QueryOrderParam{
		OrderID: orderID,
	}
}

func queryOrder(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryOrderParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeUnMarshalFailed, "failed to parse param")
	}

	okp := NewGlobalOrderKeeper(mk.marketKey, mk.cdc)
	order := okp.QueryOrder(ctx, param.OrderID)
	if order == nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeInvalidOrderID, "may be the order have deleted or not exist")
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, *order)
	if err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "could not marshal result to JSON")
	}

	return bz, nil
}

type QueryUserOrderList struct {
	User string
}

func queryUserOrderList(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryUserOrderList
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeUnMarshalFailed, "failed to parse param")
	}

	okp := NewGlobalOrderKeeper(mk.marketKey, mk.cdc)
	orders := okp.GetOrdersFromUser(ctx, param.User)

	bz, err := codec.MarshalJSONIndent(mk.cdc, orders)
	if err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}

type QueryCancelMarkets struct {
	Height int64
}

func queryWaitCancelMarkets(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryCancelMarkets
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeUnMarshalFailed, "failed to parse param")
	}

	dlk := NewDelistKeeper(mk.marketKey)
	markets := dlk.GetDelistSymbolsAtHeight(ctx, param.Height)
	bz, err := codec.MarshalJSONIndent(mk.cdc, markets)
	if err != nil {
		return nil, sdk.NewError(CodeSpaceMarket, CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}
