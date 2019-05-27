package market

import (
	"github.com/coinexchain/dex/modules/market/match"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderKeeper interface {
	Add(order *Order) error
	Exists(order *Order) bool
	Remove(order *Order) error
	GetMatchingCandidates() []match.OrderForTrade
	GetOlderThan(height int64) []*Order
}

func NewOrderKeeper(store *sdk.KVStore, symbol string) *OrderKeeper {
	return nil
}
