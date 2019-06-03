package market

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MarketInfo struct {
	Stock             string         `json:"stock"`
	Money             string         `json:"money"`
	Creator           sdk.AccAddress `json:"creator"`
	PricePrecision    byte           `json:"price_precision"`
	LastExecutedPrice sdk.Dec        `json:"last_executed_price"`
}

func (info MarketInfo) GetTags() sdk.Tags {
	return sdk.NewTags("stock", info.Stock,
		"money", info.Money,
		"creator", string(info.Creator),
		"price-precision", strconv.Itoa(int(info.PricePrecision)),
		"last-execute-price", info.LastExecutedPrice.String(),
	)
}
