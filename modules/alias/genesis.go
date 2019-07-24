package alias

import (
	"errors"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	AliasEntryList []keepers.AliasEntry `json:"alias_entry_list"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(AliasEntryList []keepers.AliasEntry) GenesisState {
	return GenesisState{
		AliasEntryList: AliasEntryList,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, entry := range data.AliasEntryList {
		keeper.AliasKeeper.AddAlias(ctx, entry.Alias, entry.Addr, entry.AsDefault, 0)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.AliasKeeper.GetAllAlias(ctx))
}

func (data GenesisState) Validate() error {
	for _, entry := range data.AliasEntryList {
		if !types.IsValidAlias(entry.Alias) {
			return errors.New("Invalid Alias")
		}
	}
	return nil
}
