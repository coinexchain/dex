package market

import sdk "github.com/cosmos/cosmos-sdk/types"

type MarketInfo struct {
	Stock             string
	Money             string
	Creator           sdk.AccAddress
	PricePrecision    byte
	LastExecutedPrice sdk.Dec
}

func (info *MarketInfo) GetTags() sdk.Tags {

	return sdk.NewTags("stock", info.Stock, "money", info.Money, "creator", info.Creator, "price-precision", info.PricePrecision, "last-execute-price", info.LastExecutedPrice)
}
