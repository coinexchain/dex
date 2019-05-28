package market

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params      Params       `json:"params"`
	Orders      []Order      `json:"orders"`
	MarketInfos []MarketInfo `json:"market_infos"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, orders []Order, infos []MarketInfo) GenesisState {
	return GenesisState{
		Params:      params,
		Orders:      orders,
		MarketInfos: infos,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []Order{}, []MarketInfo{})
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Orders {
		keeper.SetOrder(ctx, token)
	}

	for _, info := range data.MarketInfos {
		keeper.SetMarket(ctx, info)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParams(ctx), k.GetAllTokens(ctx), k.GetAllMarketInfos(ctx))
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) Validate() error {
	if err := data.Params.ValidateGenesis(); err != nil {
		return err
	}

	tokenSymbols := make(map[string]interface{})

	for _, order := range data.Orders {
		if _, exists := tokenSymbols[order.OrderID()]; exists {
			return errors.New("duplicate order found during asset ValidateGenesis")
		}
		tokenSymbols[order.OrderID()] = nil
	}

	infos := make(map[string]interface{})
	for _, info := range data.MarketInfos {
		symbol := info.Stock + SymbolSeparator + info.Money
		if _, exists := infos[symbol]; exists {
			return errors.New("duplicate market found during asset ValidateGenesis")
		}
		infos[symbol] = nil
	}
	return nil
}
