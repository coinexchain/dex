package keeper_test

import (
	"fmt"
	"github.com/coinexchain/dex/modules/authx"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func Test_queryBalances(t *testing.T) {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	params := types.DefaultParams()
	testApp.BankxKeeper.SetParams(ctx, params)

	querier := keeper.NewQuerier(testApp.BankxKeeper)
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, keeper.QueryBalances),
		Data: []byte{},
	}
	res, err := querier(ctx, []string{keeper.QueryBalances}, req)
	require.Error(t, err)
	require.Nil(t, res)

	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	req.Data = testApp.Cdc.MustMarshalJSON(keeper.NewQueryAddrBalances(addr))
	res, err = querier(ctx, []string{keeper.QueryBalances}, req)
	require.Error(t, err)
	require.Nil(t, res)

	acc := testApp.AccountKeeper.NewAccountWithAddress(ctx, addr)
	testApp.AccountKeeper.SetAccount(ctx, acc)
	res, err = querier(ctx, []string{keeper.QueryBalances}, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	testApp.AccountXKeeper.SetAccountX(ctx, authx.NewAccountXWithAddress(addr))
	res, err = querier(ctx, []string{keeper.QueryBalances}, req)
	require.NoError(t, err)
	require.NotNil(t, res)
}
