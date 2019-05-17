package asset

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the asset Querier
const (
	QueryToken     = "token"
	QueryTokenList = "tokenList"
)

// creates a querier for asset REST endpoints
func NewQuerier(tk TokenKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, req, tk)
		case QueryTokenList:
			return queryAllTokenList(ctx, req, tk)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

// defines the params for query: "custom/asset/token"
type QueryTokenParams struct {
	symbol string
}

func NewQueryAssetParams(s string) QueryTokenParams {
	return QueryTokenParams{
		symbol: s,
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, tk TokenKeeper) ([]byte, sdk.Error) {
	var params QueryTokenParams
	if err := tk.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token := tk.GetToken(ctx, params.symbol)
	if token == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("token %s does not exist", params.symbol))
	}

	bz, err := codec.MarshalJSONIndent(tk.cdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokenList(ctx sdk.Context, req abci.RequestQuery, tk TokenKeeper) ([]byte, sdk.Error) {

	tokenList := tk.GetAllTokens(ctx)
	bz, err := codec.MarshalJSONIndent(tk.cdc, tokenList)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil

}
