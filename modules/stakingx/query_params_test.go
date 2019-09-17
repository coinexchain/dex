package stakingx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

func TestQueryParams(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	params := staking.DefaultParams()
	paramsx := types.DefaultParams()
	testApp.StakingKeeper.SetParams(ctx, params)
	testApp.StakingXKeeper.SetParams(ctx, paramsx)

	querier := stakingx.NewQuerier(testApp.StakingXKeeper, testApp.Cdc)
	res, err := querier(ctx, []string{stakingx.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var mergedParams types.MergedParams
	testApp.Cdc.MustUnmarshalJSON(res, &mergedParams)
	require.Equal(t,
		string(testApp.Cdc.MustMarshalJSON(mergedParams)),
		string(testApp.Cdc.MustMarshalJSON(types.NewMergedParams(params, paramsx))))
}
