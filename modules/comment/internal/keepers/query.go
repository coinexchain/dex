package keepers

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/comment/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryCommentCount = "get-count"
)

// creates a querier for asset REST endpoints
func NewQuerier(mk Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abcitypes.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryCommentCount:
			return queryCommentCount(ctx, req, mk)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
		}
	}
}

type CommentCountInfo struct {
	CommentCount uint64 `json:"comment_count"`
}

func queryCommentCount(ctx sdk.Context, req abcitypes.RequestQuery, mk Keeper) ([]byte, sdk.Error) {
	count := mk.Cck.GetCommentCount(ctx)

	queryInfo := CommentCountInfo{
		CommentCount: count,
	}
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, queryInfo)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceComment, types.CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}
