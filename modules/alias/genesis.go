package alias

import (
	"errors"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	AliasInfoMap map[string]sdk.AccAddress `json:"alias_info_map"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(AliasInfoMap map[string]sdk.AccAddress) GenesisState {
	return GenesisState{
		AliasInfoMap: AliasInfoMap,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for alias, acc := range data.AliasInfoMap {
		keeper.AddAlias(ctx, alias, acc)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetAllAlias(ctx))
}

func (data GenesisState) Validate() error {
	for alias, _ := range data.AliasInfoMap {
		if !types.IsValidAlias(alias) {
			return errors.New("Invalid Alias")
		}
	}
	return nil
}
