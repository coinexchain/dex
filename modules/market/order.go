package market

import (
	"fmt"
	"github.com/coinexchain/dex/modules/bankx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	Sender         sdk.AccAddress
	Sequence       uint64
	Symbol         string
	OrderType      byte
	PricePrecision byte
	Price          sdk.Dec
	Quantity       sdk.Dec
	Side           byte
	TimeInForce    int
	Height         uint64

	// These field will change when order filled/cancel.
	LeftStock sdk.Dec
	Freeze    sdk.Dec
	DealStock sdk.Dec
	DealMoney sdk.Dec
}

func (or *Order) CheckOrderValidToMempool(keeper bankx.Keeper) bool {

	return true
}

func (or *Order) CheckOrderValidToExecute(keeper bankx.Keeper) bool {

	return true
}

func (or *Order) OrderID() string {
	return fmt.Sprintf("%s-%d", or.Sender, or.Sequence)
}

type OrderBook struct {
	orders     map[string]*Order
	bankKeeper bankx.Keeper
}

func (ob *OrderBook) InsertOrder(or *Order) bool {
	if !or.CheckOrderValidToMempool(ob.bankKeeper) {
		return false
	}

	if _, ok := ob.orders[or.OrderID()]; ok {
		return false
	}

	ob.orders[or.OrderID()] = or
	return true
}

func (ob *OrderBook) DelOrder(orderID string) bool {
	if _, ok := ob.orders[orderID]; !ok {
		return false
	}
	delete(ob.orders, orderID)
	return true
}
