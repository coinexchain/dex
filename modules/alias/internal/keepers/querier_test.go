package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

func TestQueryParams(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.AliasKeeper.SetParams(ctx, types.DefaultParams())

	querier := keepers.NewQuerier(testApp.AliasKeeper)
	_, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)
}
