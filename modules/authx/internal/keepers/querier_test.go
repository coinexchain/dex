package keepers_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/authx/internal/keepers"
	"github.com/coinexchain/dex/modules/authx/internal/types"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"
)

func Test_queryParams(t *testing.T) {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	params := auth.DefaultParams()
	paramsx := types.DefaultParams()
	testApp.AccountKeeper.SetParams(ctx, params)
	testApp.AccountXKeeper.SetParams(ctx, paramsx)

	querier := keepers.NewQuerier(testApp.AccountXKeeper)
	res, err := querier(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var mergedParams types.MergedParams
	testApp.Cdc.MustUnmarshalJSON(res, &mergedParams)
	require.Equal(t,
		string(testApp.Cdc.MustMarshalJSON(mergedParams)),
		string(testApp.Cdc.MustMarshalJSON(types.NewMergedParams(params, paramsx))))
}

func Test_queryAccount(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", authx.QuerierRoute, types.QueryAccountMix),
		Data: []byte{},
	}
	path0 := []string{authx.QueryAccountMix}
	query := keepers.NewQuerier(input.axk)

	res, err := query(input.ctx, path0, req)
	require.NotNil(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams([]byte("")))
	res, err = query(input.ctx, path0, req)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := testutil.KeyPubAddr()

	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams(addr))
	res, err = query(input.ctx, path0, req)
	require.NotNil(t, err)
	require.Nil(t, res)

	acc := input.ak.NewAccountWithAddress(input.ctx, addr)

	input.ak.SetAccount(input.ctx, acc)
	input.axk.SetAccountX(input.ctx, authx.NewAccountXWithAddress(addr))
	res, err = query(input.ctx, path0, req)
	require.Nil(t, err)
	require.NotNil(t, res)
}
