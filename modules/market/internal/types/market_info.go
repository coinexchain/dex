package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MarketInfo struct {
	Stock             string  `json:"stock"`
	Money             string  `json:"money"`
	PricePrecision    byte    `json:"price_precision"`
	LastExecutedPrice sdk.Dec `json:"last_executed_price"`
}
