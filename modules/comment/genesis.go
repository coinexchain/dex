package comment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	CommentCount map[string]uint64 `json:"comment_count"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(c map[string]uint64) GenesisState {
	return GenesisState{
		CommentCount: c,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make(map[string]uint64))
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for denorm, count := range data.CommentCount {
		keeper.SetCommentCount(ctx, denorm, count)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetAllCommentCount(ctx))
}

func (data GenesisState) Validate() error {
	return nil
}
