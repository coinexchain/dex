package types

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestOrder(t *testing.T) {
	bigInt, ok := sdk.NewIntFromString("57896044618658097711785492504343953926634992332820282019728792003956564819967")
	require.Equal(t, true, ok)
	bigDec := bigInt.ToDec()
	require.Equal(t, 40, len(bigDec.Int.Bytes()))
	//maxInt := sdk.NewInt(1)
	//for i:=0; i<250; i++ {
	//	maxInt = maxInt.MulRaw(2)
	//	fmt.Printf("%v %d \n", maxInt.BigInt().Bytes(), maxInt)
	//}
	maxDec := sdk.NewDec(1)
	for i := 0; i < 255; i++ {
		maxDec = maxDec.MulInt64(2)
		fmt.Printf("%v %d \n", maxDec.Int.Bytes(), maxDec)
	}
	require.Equal(t, 40, len(maxDec.Int.Bytes()))

	addr, failed := sdk.AccAddressFromHex("0123456789012345678901234567890123423456")
	require.Nil(t, failed)
	order := Order{
		Sender:   addr,
		Sequence: 9223372036854775818,
		Identify: 28,
	}
	require.Equal(t, addr.String()+"-2361183241434822609436", order.OrderID())
	order = Order{
		Sender:   addr,
		Sequence: 9223372036854775829,
		Identify: 255,
	}
	require.Equal(t, addr.String()+"-2361183241434822612479", order.OrderID())

	bz1 := DecToBigEndianBytes(sdk.NewDec(math.MaxInt64).MulInt64(100))
	bz2 := DecToBigEndianBytes(sdk.NewDec(-math.MaxInt64).MulInt64(100))
	require.Equal(t, bz1, bz2)
	bz2 = DecToBigEndianBytes(sdk.NewDec(math.MaxInt64).MulInt64(100).Add(sdk.NewDec(1)))
	require.Equal(t, 1, bytes.Compare(bz2, bz1))

	order.DealStock = 0
	order.FrozenFee = 10000
	order.Quantity = 100000
	require.Equal(t, int64(100), order.CalOrderFeeInt64(100))
	order.DealStock = 50000
	require.Equal(t, int64(5000), order.CalOrderFeeInt64(100))
	order.DealStock = 50009
	require.Equal(t, int64(5000), order.CalOrderFeeInt64(100))
	order.DealStock = 50010
	require.Equal(t, int64(5001), order.CalOrderFeeInt64(100))
	order.FrozenFee = MaxOrderAmount + 10
	order.DealStock = 100000
	require.Equal(t, MaxOrderAmount, order.CalOrderFeeInt64(100))

}
