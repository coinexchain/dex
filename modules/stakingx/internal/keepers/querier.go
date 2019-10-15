package keepers

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

// Query endpoints supported by the slashing querier
const (
	QueryPool       = "pool"
	QueryParameters = "parameters"
)

type BondPool struct {
	NotBondedTokens   sdk.Int `json:"not_bonded_tokens"`   // tokens which are not bonded to a validator (unbonded or unbonding)
	BondedTokens      sdk.Int `json:"bonded_tokens"`       // tokens which are currently bonded to a validator
	NonBondableTokens sdk.Int `json:"non_bondable_tokens"` // tokens which are in locked positions and non-bondable
	TotalSupply       sdk.Int `json:"total_supply"`        // total token supply
	BondRatio         sdk.Dec `json:"bonded_ratio"`        // bonded ratio
}

func NewQuerier(k Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryPool:
			return queryBondPool(ctx, cdc, k)
		case QueryParameters:
			return queryParameters(ctx, cdc, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown stakingx query endpoint")
		}
	}
}

func queryBondPool(ctx sdk.Context, cdc *codec.Codec, k Keeper) ([]byte, sdk.Error) {
	bondPool := k.CalcBondPoolStatus(ctx)

	res, err := codec.MarshalJSONIndent(cdc, bondPool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, cdc *codec.Codec, k Keeper) ([]byte, sdk.Error) {
	params := k.sk.GetParams(ctx)
	paramsx := k.GetParams(ctx)
	mergedParams := types.NewMergedParams(params, paramsx)

	res, err := codec.MarshalJSONIndent(cdc, mergedParams)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
