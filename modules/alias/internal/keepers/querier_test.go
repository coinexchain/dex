package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	"github.com/coinexchain/dex/testutil"
)

func TestQueryParams(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.AliasKeeper.SetParams(ctx, types.DefaultParams())

	querier := keepers.NewQuerier(testApp.AliasKeeper)
	res, err := querier(ctx, []string{keepers.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params types.Params
	testApp.Cdc.MustUnmarshalJSON(res, &params)
	require.True(t, params.Equal(types.DefaultParams()))
}

func TestQuery(t *testing.T) {
	testApp := app.NewTestApp()
	ctx := testApp.NewCtx()

	_, _, addr := testutil.KeyPubAddr()
	alias := "spiderman"

	testApp.AliasKeeper.SetParams(ctx, types.DefaultParams())
	testApp.AliasKeeper.AliasKeeper.AddAlias(ctx, alias, addr, true, 10)

	testQueryAddresses(t, testApp, ctx, addr, alias)
	testQueryAliases(t, testApp, ctx, addr, alias)
}

func testQueryAddresses(t *testing.T, testApp *app.TestApp, ctx sdk.Context, addr sdk.AccAddress, alias string) {
	reqParams := keepers.QueryAliasInfoParam{
		QueryOp: keepers.GetAddressFromAlias,
		Alias:   alias,
	}
	reqData := testApp.Cdc.MustMarshalJSON(reqParams)

	querier := keepers.NewQuerier(testApp.AliasKeeper)
	res, err := querier(ctx, []string{keepers.QueryAliasInfo}, abci.RequestQuery{Data: reqData})
	require.NoError(t, err)

	var addrs []string
	testApp.Cdc.MustUnmarshalJSON(res, &addrs)
	require.Equal(t, 1, len(addrs))
	require.Equal(t, addr.String(), addrs[0])
}

func testQueryAliases(t *testing.T, testApp *app.TestApp, ctx sdk.Context, addr sdk.AccAddress, alias string) {
	reqParams := keepers.QueryAliasInfoParam{
		QueryOp: keepers.ListAliasOfAccount,
		Owner:   addr,
	}
	reqData := testApp.Cdc.MustMarshalJSON(reqParams)

	querier := keepers.NewQuerier(testApp.AliasKeeper)
	res, err := querier(ctx, []string{keepers.QueryAliasInfo}, abci.RequestQuery{Data: reqData})
	require.NoError(t, err)

	var addrs []string
	testApp.Cdc.MustUnmarshalJSON(res, &addrs)
	require.Equal(t, 1, len(addrs))
	require.Equal(t, alias, addrs[0])
}
