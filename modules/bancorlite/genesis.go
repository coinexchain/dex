package bancorlite

import (
	"errors"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	BancorInfoMap map[string]keepers.BancorInfo `json:"bancor_info_map"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(BancorInfoMap map[string]keepers.BancorInfo) GenesisState {
	return GenesisState{
		BancorInfoMap: BancorInfoMap,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, bi := range data.BancorInfoMap {
		keeper.Bik.Save(ctx, &bi)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	m := make(map[string]keepers.BancorInfo)
	k.Bik.Iterate(ctx, func(bi *keepers.BancorInfo) {
		m[bi.Token] = *bi
	})
	return NewGenesisState(m)
}

func (data GenesisState) Validate() error {
	for token, bi := range data.BancorInfoMap {
		if token != bi.Token {
			return errors.New("Token symbol mismatch")
		}
		if token == "cet" {
			return errors.New("Token symbol can not be cet")
		}
		if len(bi.Owner) == 0 {
			return errors.New("Token owner is empty")
		}
		if len(bi.Token) == 0 {
			return errors.New("Token symbol is empty")
		}
		if !bi.MaxSupply.IsPositive() {
			return errors.New("Max Supply is not positive")
		}
		if !bi.MaxPrice.IsPositive() {
			return errors.New("Max Price is not positive")
		}
		if bi.StockInPool.IsNegative() {
			return errors.New("StockInPool is negative")
		}
		if bi.MoneyInPool.IsNegative() {
			return errors.New("MoneyInPool is negative")
		}
		if !bi.IsConsistent() {
			return errors.New("BancorInfo is not consistent")
		}
	}
	return nil
}
