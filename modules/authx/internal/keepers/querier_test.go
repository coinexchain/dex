package keepers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx/internal/keepers"
)

func Test_queryAccount(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", authx.QuerierRoute, authx.QueryAccountX),
		Data: []byte{},
	}
	path0 := []string{authx.QueryAccountX}
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

	var account authx.AccountX
	err2 := input.cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
}
