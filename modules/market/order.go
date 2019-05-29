package market

import (
	"fmt"
	"strconv"

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

	return sdk.NewTags("sender", string(or.Sender), "sequence", strconv.FormatInt(int64(or.Sequence), 10), "symbol",
		or.Symbol, "order-type", strconv.Itoa(int(or.OrderType)), "price", or.Price.String(), "quantity", strconv.FormatInt(or.Quantity, 10),
		"side", strconv.Itoa(int(or.Side)), "time-in-force", strconv.Itoa(or.TimeInForce), "height", strconv.FormatInt(or.Height, 10))
}

func (or *Order) GetTagsInOrderFilled() sdk.Tags {

	tags := or.GetTagsInOrderCreate()
	return tags.AppendTags(sdk.NewTags("left-stock", strconv.FormatInt(or.LeftStock, 10), "freeze",
		strconv.FormatInt(or.Freeze, 10), "deal-stock", strconv.FormatInt(or.DealStock, 10),
		"deal-money", strconv.FormatInt(or.DealMoney, 10)))
}
