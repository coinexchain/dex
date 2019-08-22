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
		case types.QueryAccountX:
			return queryAccountX(ctx, req, keeper)
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
