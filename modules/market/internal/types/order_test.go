package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestOrder(t *testing.T) {
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
}
