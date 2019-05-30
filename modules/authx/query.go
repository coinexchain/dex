package authx

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccountx = "accountx"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper AccountXKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAccountx:
			return queryAccountx(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown authx query endpoint")
		}
	}
}

// defines the params for query: "custom/accx/accountx"
type QueryAccountxParams struct {
	Address sdk.AccAddress
}

func NewQueryAccountxParams(addr sdk.AccAddress) QueryAccountxParams {
	return QueryAccountxParams{
		Address: addr,
	}
}

func queryAccountx(ctx sdk.Context, req abci.RequestQuery, keeper AccountXKeeper) ([]byte, sdk.Error) {
	var params QueryAccountxParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	aux, ok := keeper.GetAccountX(ctx, params.Address)
	if ok == false {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("accountx %s does not exist", params.Address))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, aux)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
