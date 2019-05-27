package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

type MarketInfo struct {
	Stock             string
	Money             string
	Creator           sdk.AccAddress
	PricePrecision    byte
	LastExecutedPrice sdk.Dec
}

func (info *MarketInfo) GetTags() sdk.Tags {

	return sdk.NewTags("stock", info.Stock, "money", info.Money, "creator", string(info.Creator), "price-precision",
		strconv.Itoa(int(info.PricePrecision)), "last-execute-price", info.LastExecutedPrice.String())
}
