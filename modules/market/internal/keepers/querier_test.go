package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/testutil"
)

func TestQueryParams(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.MarketKeeper().SetParams(ctx, types.DefaultParams())

	querier := keepers.NewQuerier(testApp.MarketKeeper())
	res, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params types.Params
	testApp.Cdc().MustUnmarshalJSON(res, &params)
	require.True(t, params.Equal(types.DefaultParams()))
}

func TestQueryMarketInfo(t *testing.T) {
	// setup
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.MarketKeeper().SetParams(ctx, types.DefaultParams())
	createMarket(ctx, testApp, "foo", "bar", 8, sdk.NewDec(10))

	// query params
	reqParams := keepers.QueryMarketParam{
		TradingPair: "foo/bar",
	}
	reqBytes := testApp.Cdc().MustMarshalJSON(reqParams)

	// query result
	querier := keepers.NewQuerier(testApp.MarketKeeper())
	resBytes, err := querier(ctx, []string{keepers.QueryMarket}, abci.RequestQuery{Data: reqBytes})
	require.NoError(t, err)
	require.NotNil(t, resBytes)

	// return data
	var res keepers.QueryMarketInfo
	testApp.Cdc().MustUnmarshalJSON(resBytes, &res)
	require.Equal(t, "foo", res.Stock)
	require.Equal(t, "bar", res.Money)
}

func createMarket(ctx sdk.Context, testApp *app.TestApp,
	stock, money string, prec byte, lep sdk.Dec) {

	_ = testApp.AssetKeeper().SetToken(ctx, &asset.BaseToken{Name: stock, Symbol: stock})
	_ = testApp.AssetKeeper().SetToken(ctx, &asset.BaseToken{Name: money, Symbol: money})
	_ = testApp.MarketKeeper().SetMarket(ctx, types.MarketInfo{
		Stock:             stock,
		Money:             money,
		PricePrecision:    prec,
		LastExecutedPrice: lep,
	})
}

func TestQueryOrder(t *testing.T) {
	// setup
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.MarketKeeper().SetParams(ctx, types.DefaultParams())
	_, _, addr := testutil.KeyPubAddr()
	order := createOrder(ctx, testApp, addr, 12345, 8)

	// query params
	reqParams := keepers.QueryOrderParam{
		OrderID: order.OrderID(),
	}
	reqBytes := testApp.Cdc().MustMarshalJSON(reqParams)

	// query result
	querier := keepers.NewQuerier(testApp.MarketKeeper())
	resBytes, err := querier(ctx, []string{keepers.QueryOrder}, abci.RequestQuery{Data: reqBytes})
	require.NoError(t, err)
	require.NotNil(t, resBytes)

	// return data
	var res types.Order
	testApp.Cdc().MustUnmarshalJSON(resBytes, &res)
	require.Equal(t, order.OrderID(), res.OrderID())
}

func createOrder(ctx sdk.Context, testApp *app.TestApp,
	sender sdk.AccAddress, seq uint64, id byte) types.Order {

	order := types.Order{
		Sender:   sender,
		Sequence: seq,
		Identify: id,
	}
	testApp.MarketKeeper().SetOrder(ctx, &order)
	return order
}

func TestQueryOrderList(t *testing.T) {
	// setup
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.MarketKeeper().SetParams(ctx, types.DefaultParams())
	_, _, addr := testutil.KeyPubAddr()
	order1 := createOrder(ctx, testApp, addr, 12345, 8)
	order2 := createOrder(ctx, testApp, addr, 12345, 9)
	order3 := createOrder(ctx, testApp, addr, 12346, 1)

	// query params
	reqParams := keepers.QueryUserOrderList{
		User: addr.String(),
	}
	reqBytes := testApp.Cdc().MustMarshalJSON(reqParams)

	// query result
	querier := keepers.NewQuerier(testApp.MarketKeeper())
	resBytes, err := querier(ctx, []string{keepers.QueryUserOrders}, abci.RequestQuery{Data: reqBytes})
	require.NoError(t, err)
	require.NotNil(t, resBytes)

	// return data
	var res []string
	testApp.Cdc().MustUnmarshalJSON(resBytes, &res)
	require.Equal(t, 3, len(res))
	require.Equal(t, order1.OrderID(), res[0])
	require.Equal(t, order2.OrderID(), res[1])
	require.Equal(t, order3.OrderID(), res[2])
}

func TestQueryWaitCancelMarkets(t *testing.T) {
	// setup
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.MarketKeeper().SetParams(ctx, types.DefaultParams())
	dlk := keepers.NewDelistKeeper(testApp.MarketKeeper().GetMarketKey())
	dlk.AddDelistRequest(ctx, 10000000, "foo/bar")

	// query params
	reqParams := keepers.QueryCancelMarkets{
		Time: 10000000,
	}
	reqBytes := testApp.Cdc().MustMarshalJSON(reqParams)

	// query result
	querier := keepers.NewQuerier(testApp.MarketKeeper())
	resBytes, err := querier(ctx, []string{keepers.QueryWaitCancelMarkets}, abci.RequestQuery{Data: reqBytes})
	require.NoError(t, err)
	require.NotNil(t, resBytes)

	// return data
	var res []string
	testApp.Cdc().MustUnmarshalJSON(resBytes, &res)
	require.Equal(t, 1, len(res))
	require.Equal(t, "foo/bar", res[0])
}
