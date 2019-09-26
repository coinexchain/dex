package keeper

import (
	"fmt"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

const (
	QueryParameters = "parameters"
	QueryBalances   = "balances"
)

// creates a querier for asset REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryParameters:
			return queryParameters(ctx, keeper)
		case QueryBalances:
			return queryBalances(ctx, keeper, req)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
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

func queryBalances(ctx sdk.Context, k Keeper, req abci.RequestQuery) ([]byte, sdk.Error) {
	var params QueryAddrBalances
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	all := struct {
		C sdk.Coins         `json:"coins"`
		L authx.LockedCoins `json:"locked_coins"`
	}{sdk.Coins{}, authx.LockedCoins{}}

	acc := params.Addr
	if au := k.ak.GetAccount(ctx, acc); au == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", acc))
	}
	all.C = k.bk.GetCoins(ctx, acc)

	if aux, ok := k.axk.GetAccountX(ctx, acc); ok {
		all.L = aux.GetAllLockedCoins()
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, all)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

type QueryAddrBalances struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAddrBalances(addr sdk.AccAddress) QueryAddrBalances {
	return QueryAddrBalances{
		Addr: addr,
	}
}
