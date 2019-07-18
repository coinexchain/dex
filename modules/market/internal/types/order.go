package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	Sender      sdk.AccAddress `json:"sender"`
	Sequence    uint64         `json:"sequence"`
	TradingPair string         `json:"trading_pair"`
	OrderType   byte           `json:"order_type"`
	Price       sdk.Dec        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Side        byte           `json:"side"`
	TimeInForce int            `json:"time_in_force"`
	Height      int64          `json:"height"`
	FrozenFee   int64          `json:"frozen_fee"`
	ExistBlocks int            `json:"exist_blocks"`

	// These fields will change when order was filled/canceled.
	LeftStock int64 `json:"left_stock"`
	Freeze    int64 `json:"freeze"`
	DealStock int64 `json:"deal_stock"`
	DealMoney int64 `json:"deal_money"`
}

func (or *Order) OrderID() string {
	// TODO. will remove the third param, ChainIDVersion
	return fmt.Sprintf("%s-%d", or.Sender, or.Sequence)
}

func (or *Order) CalOrderFee(feeForZeroDeal int64) sdk.Dec {
	actualFee := sdk.NewDec(or.DealStock).Mul(sdk.NewDec(or.FrozenFee)).Quo(sdk.NewDec(or.Quantity))
	if or.DealStock == 0 {
		actualFee = sdk.NewDec(feeForZeroDeal)
	}
	return actualFee
}
