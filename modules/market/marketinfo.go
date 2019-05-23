package market

import sdk "github.com/cosmos/cosmos-sdk/types"

type MarketInfo struct {
	Stock             string
	Money             string
	Create            sdk.AccAddress
	PricePrecision    byte
	LastExecutedPrice sdk.Dec
}

func (info *MarketInfo) CheckCreateMarketInfoValid() bool {
	return true
}
