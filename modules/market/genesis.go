package market

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

type GenesisState struct {
	Params         types.Params       `json:"params"`
	Orders         []*types.Order     `json:"orders"`
	MarketInfos    []types.MarketInfo `json:"market_infos"`
	OrderCleanTime int64              `json:"order_clean_time"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params types.Params, orders []*types.Order, infos []types.MarketInfo, cleanTime int64) GenesisState {
	return GenesisState{
		Params:         params,
		Orders:         orders,
		MarketInfos:    infos,
		OrderCleanTime: cleanTime,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams(), []*types.Order{}, []types.MarketInfo{}, 0)
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper keepers.Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Orders {
		keeper.SetOrder(ctx, token)
	}

	for _, info := range data.MarketInfos {
		keeper.SetMarket(ctx, info)
	}
	keeper.SetOrderCleanTime(ctx, data.OrderCleanTime)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, k keepers.Keeper) GenesisState {
	return NewGenesisState(k.GetParams(ctx), k.GetAllOrders(ctx), k.GetAllMarketInfos(ctx), k.GetOrderCleanTime(ctx))
}

// ValidateGenesis performs basic validation of market genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) Validate() error {
	if err := data.Params.ValidateGenesis(); err != nil {
		return err
	}

	tokenSymbols := make(map[string]struct{})

	for _, order := range data.Orders {
		if _, exists := tokenSymbols[order.OrderID()]; exists {
			return errors.New("duplicate order found during market ValidateGenesis")
		}
		tokenSymbols[order.OrderID()] = struct{}{}
	}

	infos := make(map[string]struct{})
	for _, info := range data.MarketInfos {
		symbol := info.GetSymbol()
		if _, exists := infos[symbol]; exists {
			return errors.New("duplicate market found during market ValidateGenesis")
		}
		infos[symbol] = struct{}{}
	}
	return nil
}
