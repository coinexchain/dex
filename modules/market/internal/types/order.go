package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/types"
)

type Order struct {
	Sender           sdk.AccAddress `json:"sender"`
	Sequence         uint64         `json:"sequence"`
	Identify         byte           `json:"identify"`
	TradingPair      string         `json:"trading_pair"`
	OrderType        byte           `json:"order_type"`
	Price            sdk.Dec        `json:"price"`
	Quantity         int64          `json:"quantity"`
	Side             byte           `json:"side"`
	TimeInForce      int64          `json:"time_in_force"`
	Height           int64          `json:"height"`
	FrozenCommission int64          `json:"frozen_commission"`
	ExistBlocks      int64          `json:"exist_blocks"`
	FrozenFeatureFee int64          `json:"frozen_feature_fee"`
	FrozenFee        int64          `json:"frozen_fee,omitempty"`

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

func (or *Order) CalActualOrderCommissionInt64(feeForZeroDeal int64) int64 {
	actualFee := sdk.NewDec(feeForZeroDeal)
	if or.DealStock != 0 {
		actualFee = sdk.NewDec(or.DealStock).Mul(sdk.NewDec(or.FrozenCommission)).Quo(sdk.NewDec(or.Quantity))
	}
	moa := sdk.NewDec(MaxOrderAmount)
	if actualFee.GT(moa) {
		//should not reach this clause in production, add it for safety
		actualFee = moa
	}
	return actualFee.TruncateInt64()
}

func (or *Order) CalActualOrderFeatureFeeInt64(ctx sdk.Context, freeTimeBlocks int64) int64 {
	existTime := ctx.BlockHeight() - or.Height + 1
	if existTime < freeTimeBlocks {
		return 0
	}
	chargeBlocks := existTime - freeTimeBlocks
	fee := sdk.NewDec(chargeBlocks).MulInt64(or.FrozenFeatureFee).QuoInt64(or.ExistBlocks - freeTimeBlocks).TruncateInt64()
	if fee > or.FrozenFeatureFee {
		fee = or.FrozenFeatureFee
	}
	return fee
}

func (or *Order) GetOrderUsedDenom() string {
	frozenToken, money := types.SplitSymbol(or.TradingPair)
	if or.Side == BUY {
		frozenToken = money
	}
	return frozenToken
}

func AssemblyOrderID(userAddr string, seq uint64, identify byte) string {
	idI64 := int64(identify) + 256*int64(seq%2)
	seqI64 := int64(seq / 2)
	return userAddr + OrderIDSeparator + sdk.NewInt(seqI64).MulRaw(512).AddRaw(idI64).String()
}

func DecToBigEndianBytes(d sdk.Dec) []byte {
	var result [DecByteCount]byte
	bytes := d.Int.Bytes() // returns the absolute value of d as a big-endian byte slice. sign is ignored
	count := len(bytes)
	if len(bytes) > DecByteCount {
		// Impossible to panic for cosmos 0.37.4 and golang 1.13
		// panic("dec length larger than 40")
		count = DecByteCount
	}
	for i := 1; i <= count; i++ {
		result[DecByteCount-i] = bytes[len(bytes)-i]
	}
	return result[:]
}
