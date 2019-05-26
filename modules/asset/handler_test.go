package asset

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(TokenKeeper{})

	res := h(sdk.Context{}, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized asset Msg type: "))
}

func TestIssueTokenMsg(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryToken),
		Data: []byte{},
	}
	// no token in store
	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("abc"))
	bz, err := queryToken(input.ctx, req, input.tk)
	require.Nil(t, bz)
	require.Error(t, err)

	h := NewHandler(input.tk)
	input.tk.SetParams(input.ctx, DefaultParams())
	msg := NewMsgIssueToken("ABC Token", "abc", 210000000000, tAccAddr,
		false, false, false, false)

	//case 1: issue token need valid address
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	//case2: base-case is ok
	acc := input.tk.ak.NewAccountWithAddress(input.ctx, tAccAddr)
	require.NoError(t, acc.SetCoins(types.NewCetCoins(1E18)))
	input.tk.ak.SetAccount(input.ctx, acc)

	res = h(input.ctx, msg)
	require.True(t, res.IsOK())
	bz, err = queryToken(input.ctx, req, input.tk)
	var token Token
	input.cdc.MustUnmarshalJSON(bz, &token)
	require.NoError(t, err)
	require.Equal(t, "abc", token.GetSymbol())

	// get account abc token amount
	amt := input.tk.ak.GetAccount(input.ctx, tAccAddr).GetCoins().AmountOf("abc").String()
	require.Equal(t, "210000000000", amt)
}
