package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/testapp"
)

func TestRemoveDeListRequestsBeforeTime(t *testing.T) {
	app := testapp.NewTestApp()
	ctx := app.NewCtx()
	keeper := keepers.NewDelistKeeper(app.MarketKeeper.GetMarketKey())
	keeper.AddDelistRequest(ctx, 100, "aaa/b")
	keeper.AddDelistRequest(ctx, 200, "bbb/b")
	keeper.AddDelistRequest(ctx, 300, "ccc/b")
	s := keeper.GetDelistSymbolsBeforeTime(ctx, 200)
	require.Equal(t, len(s), 2)
	require.Equal(t, s[0], "aaa/b")
	require.Equal(t, s[1], "bbb/b")
	keeper.RemoveDelistRequestsBeforeTime(ctx, 200)
	s = keeper.GetDelistSymbolsBeforeTime(ctx, 300)
	require.Equal(t, len(s), 1)
	require.Equal(t, s[0], "ccc/b")
}
