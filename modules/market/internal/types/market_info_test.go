package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarketInfo(t *testing.T) {
	require.EqualValues(t, int64(1), GetGranularityOfOrder(0))
	require.EqualValues(t, int64(10), GetGranularityOfOrder(1))
	require.EqualValues(t, int64(100), GetGranularityOfOrder(2))
	require.EqualValues(t, int64(1000), GetGranularityOfOrder(3))
	require.EqualValues(t, int64(10000), GetGranularityOfOrder(4))
	require.EqualValues(t, int64(100000), GetGranularityOfOrder(5))
	require.EqualValues(t, int64(1000000), GetGranularityOfOrder(6))
	require.EqualValues(t, int64(10000000), GetGranularityOfOrder(7))
	require.EqualValues(t, int64(100000000), GetGranularityOfOrder(8))
	require.EqualValues(t, int64(1), GetGranularityOfOrder(9))

	msg := MarketInfo{
		Stock: "abc",
		Money: "cet",
	}
	require.EqualValues(t, "abc/cet", msg.GetSymbol())
}
