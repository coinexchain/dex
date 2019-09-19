package alias

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

type GenesisState struct {
	Params         types.Params         `json:"params"`
	AliasEntryList []keepers.AliasEntry `json:"alias_entry_list"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params types.Params, AliasEntryList []keepers.AliasEntry) GenesisState {
	return GenesisState{
		Params:         params,
		AliasEntryList: AliasEntryList,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{})
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	for _, entry := range data.AliasEntryList {
		keeper.AliasKeeper.AddAlias(ctx, entry.Alias, entry.Addr, entry.AsDefault, 0)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParams(ctx), k.AliasKeeper.GetAllAlias(ctx))
}

func (data GenesisState) Validate() error {
	if err := data.Params.ValidateGenesis(); err != nil {
		return err
	}

	for _, entry := range data.AliasEntryList {
		if !types.IsValidAlias(entry.Alias) {
			return errors.New("Invalid Alias")
		}
	}
	return nil
}
