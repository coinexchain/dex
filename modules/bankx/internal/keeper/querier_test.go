package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/testapp"
)

func Test_queryParams(t *testing.T) {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	params := types.DefaultParams()
	testApp.BankxKeeper.SetParams(ctx, params)

	querier := keeper.NewQuerier(testApp.BankxKeeper)
	res, err := querier(ctx, []string{keeper.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params2 types.Params
	testApp.Cdc.MustUnmarshalJSON(res, &params2)
	require.Equal(t, params, params2)
}
