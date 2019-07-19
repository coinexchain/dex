package types

import (
	"fmt"
	"sort"

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

type PricePoint struct {
	Price     sdk.Dec `json:"price"`
	LeftStock sdk.Int `json:"left_stock"`
}

type DepthGraph struct {
	Bids []*PricePoint `json:"bids"`
	Asks []*PricePoint `json:"asks"`
}

func CalDepthGraph(orderList []*Order) *DepthGraph {
	bidMap := make(map[string]*PricePoint)
	askMap := make(map[string]*PricePoint)
	for _, order := range orderList {
		if order.Side == BID {
			p := order.Price.String()
			if _, ok := bidMap[p]; ok {
				bidMap[p].LeftStock = bidMap[p].LeftStock.AddRaw(order.LeftStock)
			} else {
				bidMap[p] = &PricePoint{Price: order.Price, LeftStock: sdk.NewInt(order.LeftStock)}
			}
		} else {
			p := order.Price.String()
			if _, ok := askMap[p]; ok {
				askMap[p].LeftStock = askMap[p].LeftStock.AddRaw(order.LeftStock)
			} else {
				askMap[p] = &PricePoint{Price: order.Price, LeftStock: sdk.NewInt(order.LeftStock)}
			}
		}
	}
	dg := &DepthGraph{
		Bids: make([]*PricePoint, 0, len(bidMap)),
		Asks: make([]*PricePoint, 0, len(askMap)),
	}
	for _, pp := range bidMap {
		dg.Bids = append(dg.Bids, pp)
	}
	for _, pp := range askMap {
		dg.Asks = append(dg.Asks, pp)
	}
	sort.Slice(dg.Bids, func(i, j int) bool {
		return dg.Bids[i].Price.GT(dg.Bids[j].Price)
	})
	sort.Slice(dg.Asks, func(i, j int) bool {
		return dg.Asks[i].Price.LT(dg.Asks[j].Price)
	})
	return dg
}
