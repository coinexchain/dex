package asset

import (
	"errors"
	"github.com/coinexchain/dex/modules/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		if err := keeper.SetToken(ctx, token); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.Whitelist {
		if err := keeper.ImportGenesisAddrKeys(ctx, types.WhitelistKey, addr); err != nil {
			panic(err)
		}
	}
	for _, addr := range data.ForbiddenAddresses {
		if err := keeper.ImportGenesisAddrKeys(ctx, types.ForbiddenAddrKey, addr); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(
		keeper.GetParams(ctx),
		keeper.GetAllTokens(ctx),
		keeper.ExportGenesisAddrKeys(ctx, types.WhitelistKey),
		keeper.ExportGenesisAddrKeys(ctx, types.ForbiddenAddrKey))
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.ValidateGenesis(); err != nil {
		return err
	}

	for _, token := range data.Tokens {
		if err := token.Validate(); err != nil {
			return err
		}
	}

	tokenSymbols := make(map[string]interface{})
	for _, token := range data.Tokens {
		if _, exists := tokenSymbols[token.GetSymbol()]; exists {
			return errors.New("duplicate token symbol found in GenesisState")
		}

		tokenSymbols[token.GetSymbol()] = nil
	}

	return nil
}
