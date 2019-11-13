package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"
)

func TestQueryParams(t *testing.T) {
	codec.RunInitFuncList()
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.BancorKeeper.SetParams(ctx, types.DefaultParams())

	querier := keepers.NewQuerier(testApp.BancorKeeper)
	res, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params types.Params
	testApp.Cdc.MustUnmarshalJSON(res, &params)
	require.True(t, params.Equal(types.DefaultParams()))
}

func TestQueryBancorInfo(t *testing.T) {
	codec.RunInitFuncList()
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()

	_, _, addr := testutil.KeyPubAddr()
	bi := keepers.BancorInfo{
		Owner:              addr,
		Stock:              "foo",
		Money:              "bar",
		InitPrice:          sdk.NewDec(10),
		MaxSupply:          sdk.NewInt(1e10),
		MaxPrice:           sdk.NewDec(10000),
		Price:              sdk.NewDec(10),
		StockInPool:        sdk.NewInt(10000),
		MoneyInPool:        sdk.NewInt(10000),
		EarliestCancelTime: 0,
	}
	testApp.BancorKeeper.Save(ctx, &bi)

	reqParams := keepers.QueryBancorInfoParam{
		Symbol: "foo/bar",
	}
	reqData := testApp.Cdc.MustMarshalJSON(reqParams)

	querier := keepers.NewQuerier(testApp.BancorKeeper)
	res, err := querier(ctx, []string{keepers.QueryBancorInfo}, abci.RequestQuery{Data: reqData})
	require.NoError(t, err)

	var bid keepers.BancorInfoDisplay
	testApp.Cdc.MustUnmarshalJSON(res, &bid)
	require.Equal(t, "foo", bid.Stock)
	require.Equal(t, "bar", bid.Money)
}
