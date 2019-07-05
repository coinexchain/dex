package asset

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the asset Querier
const (
	QueryToken           = "token-info"
	QueryTokenList       = "token-list"
	QueryWhitelist       = "token-whitelist"
	QueryForbiddenAddr   = "addr-forbidden"
	QueryReservedSymbols = "reserved-symbols"
)

// NewQuerier - creates a querier for asset REST endpoints
func NewQuerier(keeper ViewKeeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, req, keeper)
		case QueryTokenList:
			return queryAllTokenList(ctx, req, keeper)
		case QueryWhitelist:
			return queryWhitelist(ctx, req, keeper)
		case QueryForbiddenAddr:
			return queryForbiddenAddr(ctx, req, keeper)
		case QueryReservedSymbols:
			return queryReservedSymbols(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

// QueryTokenParams defines the params for query: "custom/asset/token-info"
type QueryTokenParams struct {
	Symbol string
}

func NewQueryAssetParams(s string) QueryTokenParams {
	return QueryTokenParams{
		Symbol: s,
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	var params QueryTokenParams
	if err := msgCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token := keeper.GetToken(ctx, params.Symbol)
	if token == nil {
		return nil, ErrorTokenNotFound(fmt.Sprintf("token %s not found", params.Symbol))
	}

	bz, err := codec.MarshalJSONIndent(msgCdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokenList(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(msgCdc, keeper.GetAllTokens(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// QueryWhitelistParams defines the params for query: "custom/asset/token-whitelist"
type QueryWhitelistParams struct {
	Symbol string
}

func NewQueryWhitelistParams(s string) QueryWhitelistParams {
	return QueryWhitelistParams{
		Symbol: s,
	}
}

func queryWhitelist(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	var params QueryWhitelistParams
	if err := msgCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(msgCdc, keeper.GetWhitelist(ctx, params.Symbol))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// QueryForbiddenAddrParams defines the params for query: "custom/asset/addr-forbidden"
type QueryForbiddenAddrParams struct {
	Symbol string
}

func NewQueryForbiddenAddrParams(s string) QueryForbiddenAddrParams {
	return QueryForbiddenAddrParams{
		Symbol: s,
	}
}

func queryForbiddenAddr(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	var params QueryForbiddenAddrParams
	if err := msgCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(msgCdc, keeper.GetForbiddenAddresses(ctx, params.Symbol))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryReservedSymbols(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(msgCdc, keeper.GetReservedSymbols())
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
