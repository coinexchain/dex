package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgCreateOrder(t *testing.T) {

	// Invalid address
	msg := MsgCreateOrder{}
	err := msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	msg.Sender = []byte("nihao")
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)

	// Invalid trading-pair
	msg.Sender = addr
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidSymbol, err.Code())

	msg.TradingPair = "2chs/cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidSymbol, err.Code())

	// Invalid OrderType
	msg.TradingPair = "chs/cet"
	msg.OrderType = LimitOrder + 1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidOrderType, err.Code())

	// Invalid price-precision
	msg.OrderType = LimitOrder
	msg.PricePrecision = MaxTokenPricePrecision + 1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidPricePrecision, err.Code())

	// Invalid price
	msg.PricePrecision = MaxTokenPricePrecision
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidPrice, err.Code())

	msg.Price = -1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidPrice, err.Code())

	// Invalid quantity
	msg.Price = 10
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidOrderAmount, err.Code())

	// Invalid order side
	msg.Quantity = 100
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidTradeSide, err.Code())

	// Invalid time in force
	msg.Side = BUY
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidTimeInForce, err.Code())

	// Invalid exist block
	msg.TimeInForce = GTE
	msg.ExistBlocks = -1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidExistBlocks, err.Code())

	// Success
	msg.ExistBlocks = 10000
	err = msg.ValidateBasic()
	require.EqualValues(t, nil, err)
}

func TestMsgCancelOrder(t *testing.T) {

	// Invalid address
	msg := MsgCancelOrder{}
	err := msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	msg.Sender = []byte("nihao")
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)

	// Invalid OrderID
	msg.Sender = addr
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidOrderID, err.Code())

	msg.OrderID = addr.String() + "-1-1"
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidOrderID, err.Code())

	msg.OrderID = addr.String() + "-abc"
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidOrderID, err.Code())

	msg.OrderID = "fuck-1"
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	// Success
	msg.OrderID = addr.String() + "-1"
	err = msg.ValidateBasic()
	require.EqualValues(t, nil, err)
}

func TestMsgCreateTradingPair(t *testing.T) {

	// Invalid address
	msg := MsgCreateTradingPair{}
	err := msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	msg.Creator = []byte("hello")
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidAddress, err.Code())

	// Invalid trading pair
	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)
	msg.Creator = addr
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidSymbol, err.Code())

	msg.Stock = "cet"
	msg.Money = "cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidSymbol, err.Code())

	// Invalid price precision
	msg.Stock = "chs"
	msg.Money = "cet"
	msg.PricePrecision = MaxTokenPricePrecision + 1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidPricePrecision, err.Code())

	// Success
	msg.PricePrecision = MaxTokenPricePrecision - 1
	err = msg.ValidateBasic()
	require.EqualValues(t, nil, err)
}

func TestMsgCancelTradingPair(t *testing.T) {
	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)
	msg := MsgCancelTradingPair{
		Sender:        addr,
		TradingPair:   "abc/cet",
		EffectiveTime: 10000,
	}
	err := msg.ValidateBasic()
	require.EqualValues(t, nil, err)

	msg.Sender = []byte("superman")
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidAddress(), err)

	msg.Sender = addr
	msg.TradingPair = "abc/3cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "3abc/cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "abc-cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "abc/cet"
	msg.EffectiveTime = -1
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidCancelTime(), err)
}

func TestMsgCancelOrder_GetSignBytes(t *testing.T) {
	tmp := time.Now()
	msg := MsgCancelTradingPair{EffectiveTime: tmp.UnixNano()}
	t2 := time.Unix(0, tmp.UnixNano())
	require.Equal(t, t2.UnixNano(), msg.EffectiveTime)
}

func TestMsgModifyPricePrecision(t *testing.T) {
	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)
	msg := MsgModifyPricePrecision{
		Sender:         addr,
		TradingPair:    "abc/cet",
		PricePrecision: 10,
	}
	err := msg.ValidateBasic()
	require.EqualValues(t, nil, err)

	msg.Sender = []byte("superman")
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidAddress(), err)

	msg.Sender = addr
	msg.TradingPair = "abc/3cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "3abc/cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "abc-cet"
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidSymbol(), err)

	msg.TradingPair = "abc/cet"
	msg.PricePrecision = MaxTokenPricePrecision + 1
	err = msg.ValidateBasic()
	require.EqualValues(t, ErrInvalidPricePrecision(msg.PricePrecision), err)
}
