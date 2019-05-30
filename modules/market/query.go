package market

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
)

const (
	QueryMarket = "market-info"
	QueryOrder  = "order-info"
)

// creates a querier for asset REST endpoints
func NewQuerier(mk Keeper, cdc *codec.Codec) sdk.Querier {

	return func(ctx sdk.Context, path []string, req types.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryMarket:
			return queryMarket(ctx, req, mk)
		case QueryOrder:
			return queryOrder(ctx, req, mk)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
		}
	}
}

type QueryMarketParam struct {
	Symbol string
}

func queryMarket(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryMarketParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse param: %s", err))
	}

	info, err := mk.GetMarketInfo(ctx, param.Symbol)
	if err != nil {
		return nil, sdk.ErrInternal("could not get market from blockchain" + err.Error())
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, info)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result to JSON" + err.Error())
	}
	return bz, nil
}

func NewQueryMarketParam(symbol string) QueryMarketParam {
	return QueryMarketParam{
		Symbol: symbol,
	}
}

type QueryOrderParam struct {
	Symbol  string
	OrderID string
}

func NewQueryOrderParam(symbol, orderID string) QueryOrderParam {
	return QueryOrderParam{
		Symbol:  symbol,
		OrderID: orderID,
	}
}

func queryOrder(ctx sdk.Context, req types.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryOrderParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse param: %s", err))
	}

	okp := NewOrderKeeper(mk.marketKey, param.Symbol, mk.cdc)
	order := okp.QueryOrder(ctx, param.OrderID)
	bz, err := codec.MarshalJSONIndent(mk.cdc, *order)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result to JSON" + err.Error())
	}

	return bz, nil
}
