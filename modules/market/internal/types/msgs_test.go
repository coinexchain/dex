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
	require.EqualValues(t, CodeInvalidTimeInforce, err.Code())

	// Invalid exist block
	msg.TimeInForce = GTE
	msg.ExistBlocks = -1
	err = msg.ValidateBasic()
	require.EqualValues(t, CodeInvalidExistBlocks, err.Code())
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
}

func TestMsgCancelOrder_GetSignBytes(t *testing.T) {
	tmp := time.Now()
	msg := MsgCancelTradingPair{EffectiveTime: tmp.UnixNano()}
	t2 := time.Unix(0, tmp.UnixNano())
	require.Equal(t, t2.UnixNano(), msg.EffectiveTime)
}
