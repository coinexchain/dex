package asset

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/types"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(BaseKeeper{})

	res := h(sdk.Context{}, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized asset Msg type: "))
}

func TestIssueTokenMsg(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryToken),
		Data: []byte{},
	}
	// no token in store
	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams(symbol))
	bz, err := queryToken(input.ctx, req, input.tk)
	require.Nil(t, bz)
	require.Error(t, err)

	h := NewHandler(input.tk)
	input.tk.SetParams(input.ctx, DefaultParams())
	msg := NewMsgIssueToken("ABC Token", symbol, 210000000000, tAccAddr,
		false, false, false, false)

	//case 1: issue token need valid account
	res := h(input.ctx, msg)
	require.False(t, res.IsOK())

	//case 2: base-case is ok
	err = input.tk.AddToken(input.ctx, tAccAddr, types.NewCetCoins(1E18))
	require.NoError(t, err)

	res = h(input.ctx, msg)
	require.True(t, res.IsOK())
	bz, err = queryToken(input.ctx, req, input.tk)
	var token Token
	input.cdc.MustUnmarshalJSON(bz, &token)
	require.NoError(t, err)
	require.Equal(t, symbol, token.GetSymbol())

}
