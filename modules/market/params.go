package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CreateMarketFee             = 1E11 // 1000 * 10 ^8
	FilterStaleOrderInterval    = 10
	GTEOrderLifetime            = 100
	MaxExecutedPriceChangeRatio = 25
)

type ParamsOfMarket struct {
	CreateMarketFee             sdk.Coins `json:"create_market_fee"`
	FilterStaleOrderInterval    int       `json:"filter_stale_order_interval"`
	GTEOrderLifetime            int       `json:"gte_order_lifetime"`
	MaxExecutedPriceChangeRatio int       `json:"max_executed_price_change_ratio"`
}
