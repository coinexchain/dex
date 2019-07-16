package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset/types"
)

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper BaseKeeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		if err := keeper.setToken(ctx, token); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.Whitelist {
		if err := keeper.importAddrKey(ctx, WhitelistKeyPrefix, addr); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.ForbiddenAddresses {
		if err := keeper.importAddrKey(ctx, ForbiddenAddrKeyPrefix, addr); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper BaseKeeper) types.GenesisState {
	return types.NewGenesisState(
		keeper.GetParams(ctx),
		keeper.GetAllTokens(ctx),
		keeper.ExportAddrKeys(ctx, WhitelistKeyPrefix),
		keeper.ExportAddrKeys(ctx, ForbiddenAddrKeyPrefix))
}
