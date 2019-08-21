package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

func QueryStakingParams(cdc *codec.Codec, cliCtx context.CLIContext) (staking.Params, error) {
	route := fmt.Sprintf("custom/%s/%s", staking.StoreKey, staking.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return staking.Params{}, err
	}
	var params staking.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

func QueryStakingXParams(cdc *codec.Codec, cliCtx context.CLIContext) (types.Params, error) {
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, staking.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}
