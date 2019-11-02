package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	Sender      sdk.AccAddress `json:"sender"`
	Sequence    uint64         `json:"sequence"`
	Identify    byte           `json:"identify"`
	TradingPair string         `json:"trading_pair"`
	OrderType   byte           `json:"order_type"`
	Price       sdk.Dec        `json:"price"`
	Quantity    int64          `json:"quantity"`
	Side        byte           `json:"side"`
	TimeInForce int64          `json:"time_in_force"`
	Height      int64          `json:"height"`
	FrozenFee   int64          `json:"frozen_fee"`
	ExistBlocks int64          `json:"exist_blocks"`

	// These fields will change when order was filled/canceled.
	LeftStock int64 `json:"left_stock"`
	Freeze    int64 `json:"freeze"`
	DealStock int64 `json:"deal_stock"`
	DealMoney int64 `json:"deal_money"`
}

func (or *Order) OrderID() string {
	orderID := AssemblyOrderID(or.Sender.String(), or.Sequence, or.Identify)
	return orderID
}

func (or *Order) CalOrderFeeInt64(feeForZeroDeal int64) int64 {
	actualFee := sdk.NewDec(or.DealStock).Mul(sdk.NewDec(or.FrozenFee)).Quo(sdk.NewDec(or.Quantity))
	if or.DealStock == 0 {
		actualFee = sdk.NewDec(feeForZeroDeal)
	}
	moa := sdk.NewDec(MaxOrderAmount)
	if actualFee.GT(moa) {
		//should not reach this clause in production, add it for safety
		actualFee = moa
	}
	return actualFee.TruncateInt64()
}

func AssemblyOrderID(userAddr string, seq uint64, identify byte) string {
	idI64 := int64(identify) + 256*int64(seq%2)
	seqI64 := int64(seq / 2)
	return userAddr + OrderIDSeparator + sdk.NewInt(seqI64).MulRaw(512).AddRaw(idI64).String()
}

func DecToBigEndianBytes(d sdk.Dec) []byte {
	var result [DecByteCount]byte
	bytes := d.Int.Bytes() // returns the absolute value of d as a big-endian byte slice. sign is ignored
	//todo: panic_for_test
	if len(bytes) > DecByteCount {
		panic("dec length larger than 40")
	}
	for i := 1; i <= len(bytes); i++ {
		result[DecByteCount-i] = bytes[len(bytes)-i]
	}
	return result[:]
}
