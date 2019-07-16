package asset

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_queryToken(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, QueryToken),
		Data: []byte{},
	}
	path0 := []string{QueryToken}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, []string{QueryToken}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", "abc", 2100, tAccAddr,
		false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
	require.NoError(t, err)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams(""))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("www"))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("a*B12345……6789"))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAssetParams("abc"))
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var resToken types.Token
	input.cdc.MustUnmarshalJSON(res, &resToken)
	require.Equal(t, "abc", resToken.GetSymbol())

}

func Test_queryAllTokenList(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, QueryTokenList),
		Data: []byte{},
	}
	path0 := []string{QueryTokenList}
	query := NewQuerier(input.tk)

	res, err := query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	token1, err := types.NewToken("ABC Token", "abc", 2100, tAccAddr,
		false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token1)
	require.NoError(t, err)

	token2, err := types.NewToken("XYZ Token", "xyz", 2100, tAccAddr,
		false, false, false, false, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token2)
	require.NoError(t, err)

	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var tokens []types.Token
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
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, QueryWhitelist),
		Data: []byte{},
	}
	path0 := []string{QueryWhitelist}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", symbol, 2100, tAccAddr,
		false, false, false, true, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil whitelist
	req.Data = input.cdc.MustMarshalJSON(NewQueryWhitelistParams(symbol))
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	//case 2: base-case ok
	err = input.tk.addWhitelist(input.ctx, symbol, whitelist)
	require.NoError(t, err)
	_, err = query(input.ctx, path0, req)
	require.NoError(t, err)

	err = input.tk.removeWhitelist(input.ctx, symbol, whitelist)
	require.NoError(t, err)
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

}

func Test_queryForbiddenAddr(t *testing.T) {
	input := setupTestInput()
	symbol := "abc"
	mock := mockAddresses()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, QueryForbiddenAddr),
		Data: []byte{},
	}
	path0 := []string{QueryForbiddenAddr}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", symbol, 2100, tAccAddr,
		false, false, true, true, "", "")
	require.NoError(t, err)
	err = input.tk.setToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil forbidden addr
	req.Data = input.cdc.MustMarshalJSON(NewQueryForbiddenAddrParams(symbol))
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	//case 2: base-case ok
	err = input.tk.addForbiddenAddress(input.ctx, symbol, mock)
	require.NoError(t, err)
	_, err = query(input.ctx, path0, req)
	require.NoError(t, err)

	err = input.tk.removeForbiddenAddress(input.ctx, symbol, mock)
	require.NoError(t, err)
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

}

func Test_queryReservedSymbols(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, QueryReservedSymbols),
		Data: []byte{},
	}
	path0 := []string{QueryReservedSymbols}
	query := NewQuerier(input.tk)

	res, err := query(input.ctx, path0, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	bz, _ := codec.MarshalJSONIndent(types.ModuleCdc, reserved)
	require.Equal(t, bz, res)
}

func Test_queryDefault(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, "unknown"),
		Data: []byte{},
	}
	path0 := []string{"unknown"}
	query := NewQuerier(input.tk)

	res, err := query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)
}
