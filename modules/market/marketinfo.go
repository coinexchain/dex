package market

import sdk "github.com/cosmos/cosmos-sdk/types"

type MarketInfo struct {
	Stock             string
	Money             string
	Creator           sdk.AccAddress
	PricePrecision    byte
	LastExecutedPrice sdk.Dec
}
