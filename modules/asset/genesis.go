package asset

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params  `json:"params"`
	Tokens []Token `json:"tokens"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, tokens []Token) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []Token{})
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, tk TokenKeeper, data GenesisState) {
	tk.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		if err := tk.setToken(ctx, token); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, tk TokenKeeper) GenesisState {
	return NewGenesisState(tk.GetParams(ctx), tk.GetAllTokens(ctx))
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) Validate() error {
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
