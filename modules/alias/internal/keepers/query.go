package keepers

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/alias/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryAliasInfo = "alias-info"
)

// creates a querier for asset REST endpoints
func NewQuerier(keeper Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abcitypes.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryAliasInfo:
			return queryAliasInfo(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
		}
	}
}

const (
	GetAddressFromAlias = 1
	ListAliasOfAccount  = 2
)

type QueryAliasInfoParam struct {
	Owner   sdk.AccAddress `json:"owner"`
	Alias   string         `json:"alias"`
	QueryOp int32          `json:"query_op"`
}

func queryAliasInfo(ctx sdk.Context, req abcitypes.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var param QueryAliasInfoParam
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(types.CodeSpaceAlias, types.CodeUnMarshalFailed, "failed to parse param")
	}

	var res []string
	if param.QueryOp == GetAddressFromAlias {
		acc := sdk.AccAddress(keeper.GetAddressFromAlias(ctx, param.Alias))
		if len(acc) != 0 {
			res = []string{acc.String()}
		}
	} else if param.QueryOp == ListAliasOfAccount {
		res = keeper.GetAliasListOfAccount(ctx, param.Owner)
	} else {
		return nil, sdk.NewError(types.CodeSpaceAlias, types.CodeUnknowOperation, "Unknown Operation")
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, res)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceAlias, types.CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}
