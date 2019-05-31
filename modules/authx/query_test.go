package authx

import (
	"fmt"
	"github.com/coinexchain/dex/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func Test_queryAccount(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAccountX),
		Data: []byte{},
	}

	res, err := queryAccountX(input.ctx, req, input.axk)
	require.NotNil(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams([]byte("")))
	res, err = queryAccountX(input.ctx, req, input.axk)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := testutil.KeyPubAddr()

	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams(addr))
	res, err = queryAccountX(input.ctx, req, input.axk)
	require.NotNil(t, err)
	require.Nil(t, res)

	acc := input.ak.NewAccountWithAddress(input.ctx, addr)

	input.ak.SetAccount(input.ctx, acc)
	input.axk.SetAccountX(input.ctx, NewAccountXWithAddress(addr))
	res, err = queryAccountX(input.ctx, req, input.axk)
	require.Nil(t, err)
	require.NotNil(t, res)

	var account AccountX
	err2 := input.cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
}
