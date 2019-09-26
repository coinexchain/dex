package keepers

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx/internal/types"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper AccountXKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		case types.QueryAccountX:
			return queryAccountX(ctx, req, keeper)
		case types.QueryAccountMix:
			return queryAccountMix(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown authx query endpoint")
		}
	}
}

func queryAccountX(ctx sdk.Context, req abci.RequestQuery, keeper AccountXKeeper) ([]byte, sdk.Error) {
	var params auth.QueryAccountParams

	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	aux, ok := keeper.GetAccountX(ctx, params.Address)
	if !ok {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("accountx %s does not exist", params.Address))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, aux)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAccountMix(ctx sdk.Context, req abci.RequestQuery, keeper AccountXKeeper) ([]byte, sdk.Error) {
	var params auth.QueryAccountParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	addr := params.Address
	au := keeper.ak.GetAccount(ctx, addr)
	if au == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	aux, ok := keeper.GetAccountX(ctx, addr)
	if !ok {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("accountx %s does not exist", addr))
	}

	mix := types.NewAccountMix(au, aux)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, mix)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
func queryParameters(ctx sdk.Context, k AccountXKeeper) ([]byte, sdk.Error) {
	params := k.ak.GetParams(ctx)
	paramsx := k.GetParams(ctx)
	mergedParams := types.NewMergedParams(params, paramsx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, mergedParams)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
