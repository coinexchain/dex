package keepers

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/asset/internal/types"
)

// NewQuerier - creates a querier for asset REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		case types.QueryToken:
			return queryToken(ctx, req, keeper)
		case types.QueryTokenList:
			return queryAllTokenList(ctx, req, keeper)
		case types.QueryWhitelist:
			return queryWhitelist(ctx, req, keeper)
		case types.QueryForbiddenAddr:
			return queryForbiddenAddr(ctx, req, keeper)
		case types.QueryReservedSymbols:
			return queryReservedSymbols()
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper TokenKeeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token := keeper.GetToken(ctx, params.Symbol)
	if token == nil {
		return nil, types.ErrTokenNotFound(params.Symbol)
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokenList(ctx sdk.Context, req abci.RequestQuery, keeper TokenKeeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetAllTokens(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryWhitelist(ctx sdk.Context, req abci.RequestQuery, keeper TokenKeeper) ([]byte, sdk.Error) {
	var params types.QueryWhitelistParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetWhitelist(ctx, params.Symbol))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryForbiddenAddr(ctx sdk.Context, req abci.RequestQuery, keeper TokenKeeper) ([]byte, sdk.Error) {
	var params types.QueryForbiddenAddrParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, keeper.GetForbiddenAddresses(ctx, params.Symbol))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryReservedSymbols() ([]byte, sdk.Error) {
	reserved := types.GetReservedSymbols()
	var s = ""
	for _, symbol := range reserved {
		s = s + symbol + ", "
	}
	res := "[" + s + "]"
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
