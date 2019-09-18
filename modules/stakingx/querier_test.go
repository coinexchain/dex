package stakingx_test

import (
	"github.com/coinexchain/dex/modules/stakingx"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewQuerier(t *testing.T) {
	//intialize
	sxk, ctx, _ := setUpInput()
	cdc := codec.New()

	sxk.SetParams(ctx, stakingx.DefaultParams())

	//query succeed
	querier := stakingx.NewQuerier(sxk.Keeper, cdc)
	path := stakingx.QueryPool

	_, err := querier(ctx, []string{path}, abci.RequestQuery{})
	require.Nil(t, err)

	//query fail
	failPath := "fake"
	_, err = querier(ctx, []string{failPath}, abci.RequestQuery{})
	require.Equal(t, sdk.CodeUnknownRequest, err.Code())
}
