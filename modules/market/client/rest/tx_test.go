package rest

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

func TestCmd(t *testing.T) {
	createMarket := createMarketReq{
		Stock:          "etc",
		Money:          "cet",
		PricePrecision: 8,
	}
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	msg, _ := createMarket.GetMsg(nil, addr)
	assert.Equal(t, types.MsgCreateTradingPair{
		Stock:          "etc",
		Money:          "cet",
		Creator:        addr,
		PricePrecision: 8,
	}, msg)
	//==============
	cancelMarket := cancelMarketReq{
		TradingPair: "etc/cet",
		Time:        12345678,
	}
	msg, _ = cancelMarket.GetMsg(nil, addr)
	assert.Equal(t, types.MsgCancelTradingPair{
		Sender:        addr,
		TradingPair:   "etc/cet",
		EffectiveTime: 12345678,
	}, msg)
	//==============
	req := modifyPricePrecision{
		TradingPair:    "etc/cet",
		PricePrecision: 9,
	}
	msg, _ = req.GetMsg(nil, addr)
	assert.Equal(t, types.MsgModifyPricePrecision{
		Sender:         addr,
		TradingPair:    "etc/cet",
		PricePrecision: 9,
	}, msg)
	//==============
	createOrder := createOrderReq{
		OrderType:      types.LIMIT,
		TradingPair:    "etc/cet",
		Identify:       0,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		ExistBlocks:    25000,
		TimeInForce:    types.GTE,
	}
	httpReq, _ := http.NewRequest("POST", "http://example.com/market/gte-orders", nil)
	msg, _ = createOrder.GetMsg(httpReq, addr)
	assert.Equal(t, types.MsgCreateOrder{
		Sender:         addr,
		Identify:       0,
		TradingPair:    "etc/cet",
		OrderType:      types.LIMIT,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		TimeInForce:    types.GTE,
		ExistBlocks:    25000,
	}, msg)
	httpReq, _ = http.NewRequest("POST", "http://example.com/market/ioc-orders", nil)
	msg, _ = createOrder.GetMsg(httpReq, addr)
	assert.Equal(t, types.MsgCreateOrder{
		Sender:         addr,
		Identify:       0,
		TradingPair:    "etc/cet",
		OrderType:      types.LIMIT,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		TimeInForce:    types.IOC,
		ExistBlocks:    25000,
	}, msg)
	//==============
	cancelOrder := cancelOrderReq{
		OrderID: "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}
	msg, _ = cancelOrder.GetMsg(nil, addr)
	assert.Equal(t, &types.MsgCancelOrder{
		Sender:  addr,
		OrderID: "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}, msg)
}
