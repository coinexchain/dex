package authx_test

//
//import (
//	"fmt"
//	"github.com/coinexchain/dex/modules/authx"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//	abci "github.com/tendermint/tendermint/abci/types"
//
//	"github.com/coinexchain/dex/modules/authx/types"
//	"github.com/coinexchain/dex/testutil"
//	"github.com/cosmos/cosmos-sdk/x/auth"
//)
//
//func Test_queryAccount(t *testing.T) {
//	input := setupTestInput()
//	req := abci.RequestQuery{
//		Path: fmt.Sprintf("custom/%s/%s", authx.QuerierRoute, authx.QueryAccountX),
//		Data: []byte{},
//	}
//
//	res, err := queryAccountX(input.ctx, req, input.axk)
//	require.NotNil(t, err)
//	require.Nil(t, res)
//
//	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams([]byte("")))
//	res, err = queryAccountX(input.ctx, req, input.axk)
//	require.NotNil(t, err)
//	require.Nil(t, res)
//
//	_, _, addr := testutil.KeyPubAddr()
//
//	req.Data = input.cdc.MustMarshalJSON(auth.NewQueryAccountParams(addr))
//	res, err = queryAccountX(input.ctx, req, input.axk)
//	require.NotNil(t, err)
//	require.Nil(t, res)
//
//	acc := input.ak.NewAccountWithAddress(input.ctx, addr)
//
//	input.ak.SetAccount(input.ctx, acc)
//	input.axk.SetAccountX(input.ctx, types.NewAccountXWithAddress(addr))
//	res, err = queryAccountX(input.ctx, req, input.axk)
//	require.Nil(t, err)
//	require.NotNil(t, res)
//
//	var account types.AccountX
//	err2 := input.cdc.UnmarshalJSON(res, &account)
//	require.Nil(t, err2)
//}
