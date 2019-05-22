package market

import "github.com/cosmos/cosmos-sdk/types"

type MarketInfo struct {
	Stock             string
	Money             string
	Create            string
	PricePrecision    byte
	LastExecutedPrice types.Dec
}

func (info *MarketInfo) CheckCreateMarketInfoValid() bool {
	return true
}
