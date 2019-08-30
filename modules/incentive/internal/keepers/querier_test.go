package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

func TestQueryParams(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	params := types.DefaultParams()
	testApp.IncentiveKeeper.SetParams(ctx, params)

	querier := keepers.NewQuerier(testApp.IncentiveKeeper)
	res, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params2 types.Params
	testApp.Cdc.MustUnmarshalJSON(res, &params2)
	require.Equal(t, params, params2)
}
