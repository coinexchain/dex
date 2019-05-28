package asset

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_queryToken(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryToken),
		Data: []byte{},
	}

	// no token
	res, err := queryToken(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := NewToken("ABC Token", "abc", 2100, tAccAddr,
		false, false, false, false)
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token)
	require.NoError(t, err)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams(""))
	res, err = queryToken(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("www"))
	res, err = queryToken(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("a*B12345……6789"))
	res, err = queryToken(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("abc"))
	res, err = queryToken(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.NotNil(t, res)

	var resToken Token
	input.cdc.MustUnmarshalJSON(res, &resToken)
	require.Equal(t, "abc", resToken.GetSymbol())

}

func Test_queryAllTokenList(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryTokenList),
		Data: []byte{},
	}

	res, err := queryAllTokenList(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	token1, err := NewToken("ABC Token", "abc", 2100, tAccAddr,
		false, false, false, false)
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token1)
	require.NoError(t, err)

	token2, err := NewToken("XYZ Token", "xyz", 2100, tAccAddr, false, false, false, false)
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token2)
	require.NoError(t, err)

	res, err = queryAllTokenList(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.NotNil(t, res)

	var tokens []Token
	input.cdc.MustUnmarshalJSON(res, &tokens)
	require.Equal(t, 2, len(tokens))
	require.Contains(t, []string{"abc", "xyz"}, tokens[0].GetSymbol())
	require.Contains(t, []string{"abc", "xyz"}, tokens[1].GetSymbol())
}
