package keepers

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

const (
	QueryMarket            = "market-info"
	QueryMarkets           = "market-list"
	QueryOrder             = "order-info"
	QueryUserOrders        = "user-order-list"
	QueryWaitCancelMarkets = "wait-cancel-markets"
	QueryParameters        = "parameters"
)

// creates a querier for asset REST endpoints
func NewQuerier(mk Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryParameters:
			return queryParameters(ctx, mk)
		case QueryMarket:
			return queryMarket(ctx, req, mk)
		case QueryMarkets:
			return queryMarketList(ctx, req, mk)
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

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

type QueryMarketParam struct {
	TradingPair string
}

func NewQueryMarketParam(symbol string) QueryMarketParam {
	return QueryMarketParam{
		TradingPair: symbol,
	}
}

type QueryMarketInfo struct {
	Creator           sdk.AccAddress `json:"creator"`
	Stock             string         `json:"stock"`
	Money             string         `json:"money"`
	PricePrecision    string         `json:"price_precision"`
	LastExecutedPrice sdk.Dec        `json:"last_executed_price"`
}

func queryMarket(ctx sdk.Context, req abci.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryMarketParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse param: %s", err))
	}

	info, err := mk.GetMarketInfo(ctx, param.TradingPair)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeInvalidSymbol, "may be the market have deleted or not exist")
	}

	queryInfo := QueryMarketInfo{
		Creator:           mk.MarketOwner(ctx, info),
		Stock:             info.Stock,
		Money:             info.Money,
		PricePrecision:    strconv.Itoa(int(info.PricePrecision)),
		LastExecutedPrice: info.LastExecutedPrice,
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, queryInfo)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}

type MarketInfoList []QueryMarketInfo

func queryMarketList(ctx sdk.Context, req abci.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	infos := mk.GetAllMarketInfos(ctx)
	mInfoList := make([]QueryMarketInfo, len(infos))

	for i, info := range infos {
		mInfoList[i] = QueryMarketInfo{
			Creator:           mk.MarketOwner(ctx, info),
			Stock:             info.Stock,
			Money:             info.Money,
			PricePrecision:    strconv.Itoa(int(info.PricePrecision)),
			LastExecutedPrice: info.LastExecutedPrice,
		}
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, mInfoList)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeMarshalFailed, "could not marshal result to JSON")
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

func queryOrder(ctx sdk.Context, req abci.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryOrderParam
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeUnMarshalFailed, "failed to parse param")
	}

	okp := NewGlobalOrderKeeper(mk.marketKey, mk.cdc)
	order := okp.QueryOrder(ctx, param.OrderID)
	if order == nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeInvalidOrderID, "may be the order have deleted or not exist")
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, *order)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeMarshalFailed, "could not marshal result to JSON")
	}

	return bz, nil
}

type QueryUserOrderList struct {
	User string
}

func queryUserOrderList(ctx sdk.Context, req abci.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryUserOrderList
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeUnMarshalFailed, "failed to parse param")
	}

	okp := NewGlobalOrderKeeper(mk.marketKey, mk.cdc)
	orders := okp.GetOrdersFromUser(ctx, param.User)
	if len(orders) == 0 {
		orders = append(orders, "")
	}
	bz, err := codec.MarshalJSONIndent(mk.cdc, orders)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeMarshalFailed, "could not marshal result to JSON")
	}

	return bz, nil
}

type QueryCancelMarkets struct {
	Time int64
}

func queryWaitCancelMarkets(ctx sdk.Context, req abci.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	var param QueryCancelMarkets
	if err := mk.cdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeUnMarshalFailed, "failed to parse param")
	}

	dlk := NewDelistKeeper(mk.marketKey)
	markets := dlk.GetDelistSymbolsBeforeTime(ctx, param.Time)
	bz, err := codec.MarshalJSONIndent(mk.cdc, markets)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceMarket, types.CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}
