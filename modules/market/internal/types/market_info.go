package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dex "github.com/coinexchain/dex/types"
)

type MarketInfo struct {
	Stock             string  `json:"stock"`
	Money             string  `json:"money"`
	PricePrecision    byte    `json:"price_precision"`
	LastExecutedPrice sdk.Dec `json:"last_executed_price"`
	OrderPrecision    byte    `json:"order_precision"`
}

func GetGranularityOfOrder(orderPrecision byte) int64 {
	if orderPrecision == 0 {
		return 1
	}
	return int64(math.Pow10(int(8 - (orderPrecision - 1))))
}

func (msg MarketInfo) GetSymbol() string {
	return dex.GetSymbol(msg.Stock, msg.Money)
}
