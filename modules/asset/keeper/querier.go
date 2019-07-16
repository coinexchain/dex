package keeper

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)



// NewQuerier - creates a querier for asset REST endpoints
func NewQuerier(keeper ViewKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, keeper)
		case types.QueryTokenList:
			return queryAllTokenList(ctx, req, keeper)
		case types.QueryWhitelist:
			return queryWhitelist(ctx, req, keeper)
		case types.QueryForbiddenAddr:
			return queryForbiddenAddr(ctx, req, keeper)
		case types.QueryReservedSymbols:
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
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token := keeper.GetToken(ctx, params.Symbol)
	if token == nil {
		return nil, types.ErrorTokenNotFound(fmt.Sprintf("token %s not found", params.Symbol))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokenList(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetAllTokens(ctx))
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
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetWhitelist(ctx, params.Symbol))
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
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetForbiddenAddresses(ctx, params.Symbol))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryReservedSymbols(ctx sdk.Context, req abci.RequestQuery, keeper ViewKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, types.GetReservedSymbols())
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
