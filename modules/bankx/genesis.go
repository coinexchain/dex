package bankx

import (
	"github.com/coinexchain/dex/modules/authx"
	gaia_app "github.com/cosmos/cosmos-sdk/cmd/gaia/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Param Param `json:"param"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(param Param) GenesisState {
	return GenesisState{
		Param: param,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParam())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState,
	accounts []gaia_app.GenesisAccount) {

	keeper.SetParam(ctx, data.Param)
	activateGenesisAccounts(ctx, keeper, accounts)
}

func activateGenesisAccounts(ctx sdk.Context, keeper Keeper, accounts []gaia_app.GenesisAccount) {
	for _, acc := range accounts {
		accX := authx.AccountX{Address: acc.Address, Activated: true}
		keeper.axk.SetAccountX(ctx, accX)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParam(ctx)
	return NewGenesisState(params)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	activatedFee := data.Param.ActivatedFee
	if activatedFee < 0 {
		return sdk.NewError(CodeSpaceBankx, CodeInvalidActivatedFee, "invalid activated fees")
	}
	return nil
}
