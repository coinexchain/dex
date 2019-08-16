package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/incentive"
)

type GenesisState struct {
	Params Params `json:"params"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// cache non-bondable addresses
	addresses := keeper.getAllVestingAccountAddresses(ctx)
	addresses = append(addresses, incentive.PoolAddr)
	if cetOwner := keeper.getCetOwnerAddress(ctx); cetOwner != nil {
		addresses = append(addresses, cetOwner)
	}
	keeper.setNonBondableAddresses(ctx, addresses)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) ValidateGenesis() error {
	msd := data.Params.MinSelfDelegation
	if !msd.IsPositive() {
		return ErrInvalidMinSelfDelegation(msd)
	}
	return nil
}
