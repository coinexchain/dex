package asset

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the asset Querier
const (
	QueryToken     = "token-info"
	QueryTokenList = "token-list"
	QueryWhitelist = "token-whitelist"
)

// creates a querier for asset REST endpoints
func NewQuerier(tk TokenKeeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, req, tk)
		case QueryTokenList:
			return queryAllTokenList(ctx, req, tk)
		case QueryWhitelist:
			return queryWhitelist(ctx, req, tk)

		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

// defines the params for query: "custom/asset/token-info"
type QueryTokenParams struct {
	Symbol string
}

func NewQueryAssetParams(s string) QueryTokenParams {
	return QueryTokenParams{
		Symbol: s,
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, tk TokenKeeper) ([]byte, sdk.Error) {
	var params QueryTokenParams
	if err := tk.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token := tk.GetToken(ctx, params.Symbol)
	if token == nil {
		return nil, ErrorTokenNotFound(fmt.Sprintf("token %s not found", params.Symbol))
	}

	bz, err := codec.MarshalJSONIndent(tk.cdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokenList(ctx sdk.Context, req abci.RequestQuery, tk TokenKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(tk.cdc, tk.GetAllTokens(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// defines the params for query: "custom/asset/token-whitelist"
type QueryWhitelistParams struct {
	Symbol string
}

func NewQueryWhitelistParams(s string) QueryWhitelistParams {
	return QueryWhitelistParams{
		Symbol: s,
	}
}

func queryWhitelist(ctx sdk.Context, req abci.RequestQuery, tk TokenKeeper) ([]byte, sdk.Error) {
	var params QueryWhitelistParams
	if err := tk.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	whitelist := tk.GetWhitelist(ctx, params.Symbol)
	bz, err := codec.MarshalJSONIndent(tk.cdc, whitelist)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
