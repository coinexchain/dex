package market

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	Sender      sdk.AccAddress
	Sequence    uint64
	Symbol      string
	OrderType   byte
	Price       sdk.Dec
	Quantity    int64
	Side        byte
	TimeInForce int
	Height      int64

	// These field will change when order filled/cancel.
	LeftStock int64
	Freeze    int64
	DealStock int64
	DealMoney int64
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
