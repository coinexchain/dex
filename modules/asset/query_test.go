package asset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

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
		false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
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
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	token1, err := NewToken("ABC Token", "abc", 2100, tAccAddr,
		false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token1)
	require.NoError(t, err)

	token2, err := NewToken("XYZ Token", "xyz", 2100, tAccAddr, false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token2)
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

func Test_queryWhitelist(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	whitelist := mockWhitelist()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryWhitelist),
		Data: []byte{},
	}

	// no token
	res, err := queryWhitelist(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := NewToken("ABC Token", symbol, 2100, tAccAddr,
		false, false, false, true, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil whitelist
	req.Data = input.cdc.MustMarshalJSON(NewQueryWhitelistParams(symbol))
	res, err = queryWhitelist(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	//case 2: base-case ok
	err = input.tk.addWhitelist(input.ctx, symbol, whitelist)
	require.NoError(t, err)
	_, err = queryWhitelist(input.ctx, req, input.tk)
	require.NoError(t, err)

	err = input.tk.removeWhitelist(input.ctx, symbol, whitelist)
	require.NoError(t, err)
	res, err = queryWhitelist(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

}

func Test_queryForbiddenAddr(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	mock := mockAddresses()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", RouterKey, QueryForbiddenAddr),
		Data: []byte{},
	}

	// no token
	res, err := queryForbiddenAddr(input.ctx, req, input.tk)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := NewToken("ABC Token", symbol, 2100, tAccAddr,
		false, false, true, true, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil forbidden addr
	req.Data = input.cdc.MustMarshalJSON(NewQueryForbiddenAddrParams(symbol))
	res, err = queryForbiddenAddr(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	//case 2: base-case ok
	err = input.tk.addForbidAddress(input.ctx, symbol, mock)
	require.NoError(t, err)
	_, err = queryForbiddenAddr(input.ctx, req, input.tk)
	require.NoError(t, err)

	err = input.tk.removeForbidAddress(input.ctx, symbol, mock)
	require.NoError(t, err)
	res, err = queryForbiddenAddr(input.ctx, req, input.tk)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

}
