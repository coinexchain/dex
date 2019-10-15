package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/modules/stakingx/internal/keepers"
	"github.com/coinexchain/dex/modules/stakingx/internal/types"
	"github.com/coinexchain/dex/testapp"
)

func TestNewQuerier(t *testing.T) {
	//intialize
	sxk, ctx, _ := setUpInput()
	cdc := codec.New()

	sxk.SetParams(ctx, stakingx.DefaultParams())

	//query succeed
	querier := keepers.NewQuerier(sxk.Keeper, cdc)
	path := keepers.QueryPool

	_, err := querier(ctx, []string{path}, abci.RequestQuery{})
	require.Nil(t, err)

	//query fail
	failPath := "fake"
	_, err = querier(ctx, []string{failPath}, abci.RequestQuery{})
	require.Equal(t, sdk.CodeUnknownRequest, err.Code())
}

func TestQueryParams(t *testing.T) {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	params := staking.DefaultParams()
	paramsx := types.DefaultParams()
	testApp.StakingKeeper.SetParams(ctx, params)
	testApp.StakingXKeeper.SetParams(ctx, paramsx)

	querier := keepers.NewQuerier(testApp.StakingXKeeper, testApp.Cdc)
	res, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var mergedParams types.MergedParams
	testApp.Cdc.MustUnmarshalJSON(res, &mergedParams)
	require.Equal(t,
		string(testApp.Cdc.MustMarshalJSON(mergedParams)),
		string(testApp.Cdc.MustMarshalJSON(types.NewMergedParams(params, paramsx))))
}
