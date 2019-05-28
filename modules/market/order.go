package market

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	Sender      sdk.AccAddress `json:"sender"`
	Sequence    uint64         `json:"sequence"`
	Symbol      string         `json:"symbol"`
	OrderType   byte           `json:"order_type"`
	Price       sdk.Dec        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Side        byte           `json:"side"`
	TimeInForce int            `json:"time_in_force"`
	Height      int64          `json:"height"`

	// These field will change when order filled/cancel.
	LeftStock int64 `json:"left_stock"`
	Freeze    int64 `json:"freeze"`
	DealStock int64 `json:"deal_stock"`
	DealMoney int64 `json:"deal_money"`
}

func (or *Order) OrderID() string {
	return fmt.Sprintf("%s-%d", or.Sender, or.Sequence)
}

func (or *Order) GetTagsInOrderCreate() sdk.Tags {

	return sdk.NewTags("sender", or.Sender, "sequence", or.Sequence, "symbol",
		or.Symbol, "order-type", or.OrderType, "price", or.Price, "quantity", or.Quantity,
		"side", or.Side, "time-in-force", or.TimeInForce, "height", or.Height)
}

func (or *Order) GetTagsInOrderFilled() sdk.Tags {

	tags := or.GetTagsInOrderCreate()
	return tags.AppendTags(sdk.NewTags("left-stock", or.LeftStock, "freeze",
		or.Freeze, "deal-stock", or.DealStock, "deal-money", or.DealMoney))
}
