package types

import (
	"fmt"

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
	orderID, _ := AssemblyOrderID(or.Sender.String(), or.Sequence, or.Identify)
	return orderID
}

func (or *Order) CalOrderFee(feeForZeroDeal int64) sdk.Dec {
	actualFee := sdk.NewDec(or.DealStock).Mul(sdk.NewDec(or.FrozenFee)).Quo(sdk.NewDec(or.Quantity))
	if or.DealStock == 0 {
		actualFee = sdk.NewDec(feeForZeroDeal)
	}
	return actualFee.TruncateDec()
}

func AssemblyOrderID(userAddr string, seq uint64, identify byte) (string, error) {
	seqInt, ok := sdk.NewIntFromString(fmt.Sprintf("%d", seq))
	if !ok {
		return "", fmt.Errorf("invalid sequence : %d", seq)
	}
	orderID := userAddr + OrderIDSeparator + seqInt.MulRaw(256).AddRaw(int64(identify)).String()
	return orderID, nil
}

type PricePoint struct {
	Price     sdk.Dec `json:"price"`
	LeftStock sdk.Int `json:"left_stock"`
}

type DepthGraph struct {
	Bids []*PricePoint `json:"bids"`
	Asks []*PricePoint `json:"asks"`
}

func CalDepthGraph(orderList []*Order) *DepthGraph {
	bidMap := make(map[string]int)
	askMap := make(map[string]int)
	bidList := make([]*PricePoint, 0, 100)
	askList := make([]*PricePoint, 0, 100)
	for _, order := range orderList {
		p := string(DecToBigEndianBytes(order.Price))
		if order.Side == BID {
			if offset, ok := bidMap[p]; ok {
				bidList[offset].LeftStock = bidList[offset].LeftStock.AddRaw(order.LeftStock)
			} else {
				bidList = append(bidList, &PricePoint{Price: order.Price, LeftStock: sdk.NewInt(order.LeftStock)})
				bidMap[p] = len(bidList)
			}
		} else {
			if offset, ok := askMap[p]; ok {
				askList[offset].LeftStock = askList[offset].LeftStock.AddRaw(order.LeftStock)
			} else {
				askList = append(askList, &PricePoint{Price: order.Price, LeftStock: sdk.NewInt(order.LeftStock)})
				askMap[p] = len(askList)
			}
		}
	}
	return &DepthGraph{Bids: bidList, Asks: askList}
}

func DecToBigEndianBytes(d sdk.Dec) []byte {
	var result [DecByteCount]byte
	bytes := d.Int.Bytes() //  returns the absolute value of d as a big-endian byte slice.
	for i := 1; i <= len(bytes); i++ {
		result[DecByteCount-i] = bytes[len(bytes)-i]
	}
	return result[:]
}
