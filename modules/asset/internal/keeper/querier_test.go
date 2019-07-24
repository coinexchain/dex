package keeper

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_queryToken(t *testing.T) {
	input := createTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryToken),
		Data: []byte{},
	}
	path0 := []string{types.QueryToken}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, []string{types.QueryToken}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", "")
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token)
	require.NoError(t, err)

	req.Data = input.cdc.MustMarshalJSON(types.NewQueryAssetParams(""))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(types.NewQueryAssetParams("www"))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(types.NewQueryAssetParams("a*B12345……6789"))
	res, err = query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(types.NewQueryAssetParams("abc"))
	res, err = query(input.ctx, path0, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var resToken types.Token
	input.cdc.MustUnmarshalJSON(res, &resToken)
	require.Equal(t, "abc", resToken.GetSymbol())

}

func Test_queryAllTokenList(t *testing.T) {
	input := createTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryTokenList),
		Data: []byte{},
	}
	path0 := []string{types.QueryTokenList}
	query := NewQuerier(input.tk)

	res, err := query(input.ctx, path0, req)
	require.NoError(t, err)
	require.Equal(t, []byte("[]"), res)

	token1, err := types.NewToken("ABC Token", "abc", sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", "")
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token1)
	require.NoError(t, err)

	token2, err := types.NewToken("XYZ Token", "xyz", sdk.NewInt(2100), testAddr,
		false, false, false, false, "", "", "")
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token2)
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
	input := createTestInput()
	symbol := "abc"
	whitelist := mockAddrList()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryWhitelist),
		Data: []byte{},
	}
	path0 := []string{types.QueryWhitelist}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		false, false, false, true, "", "", "")
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil whitelist
	req.Data = input.cdc.MustMarshalJSON(types.NewQueryWhitelistParams(symbol))
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
	input := createTestInput()
	symbol := "abc"
	mock := mockAddrList()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryForbiddenAddr),
		Data: []byte{},
	}
	path0 := []string{types.QueryForbiddenAddr}
	query := NewQuerier(input.tk)

	// no token
	res, err := query(input.ctx, path0, req)
	require.Error(t, err)
	require.Nil(t, res)

	// set token
	token, err := types.NewToken("ABC Token", symbol, sdk.NewInt(2100), testAddr,
		false, false, true, true, "", "", "")
	require.NoError(t, err)
	err = input.tk.SetToken(input.ctx, token)
	require.NoError(t, err)

	//case 1: nil forbidden addr
	req.Data = input.cdc.MustMarshalJSON(types.NewQueryForbiddenAddrParams(symbol))
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
	input := createTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryReservedSymbols),
		Data: []byte{},
	}
	path0 := []string{types.QueryReservedSymbols}
	query := NewQuerier(input.tk)

	res, err := query(input.ctx, path0, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	bz, _ := codec.MarshalJSONIndent(types.ModuleCdc, types.GetReservedSymbols())
	require.Equal(t, bz, res)
}

func Test_queryDefault(t *testing.T) {
	input := createTestInput()
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
