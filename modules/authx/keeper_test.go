package authx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetSetParams(t *testing.T) {
	input := setupTestInput()
	params := DefaultParams()
	input.axk.SetParams(input.ctx, params)
	params2 := input.axk.GetParams(input.ctx)
	require.True(t, params.Equal(params2))
}

func TestAccountXGetSet(t *testing.T) {
	input := setupTestInput()
	addr := sdk.AccAddress([]byte("some-address"))

	_, ok := input.axk.GetAccountX(input.ctx, addr)
	require.False(t, ok)

	//create account
	acc := NewAccountXWithAddress(addr)
	require.Equal(t, addr, acc.Address)

	input.axk.SetAccountX(input.ctx, acc)

	acc, ok = input.axk.GetAccountX(input.ctx, addr)
	require.True(t, ok)

	acc.MemoRequired = false
	input.axk.SetAccountX(input.ctx, acc)
	acc, _ = input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t, false, acc.MemoRequired)

	lockedcoin := acc.LockedCoins
	require.Nil(t, lockedcoin)
}

func TestAddressStoreKey(t *testing.T) {
	addr := sdk.AccAddress([]byte("some-address1"))
	addrStoreKey := AddressStoreKey(addr)
	expectedOutput := []byte{0x1, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x31}
	require.Equal(t, expectedOutput, addrStoreKey)
}

func TestGetOrCreateAccountX(t *testing.T) {
	input := setupTestInput()
	addr := sdk.AccAddress([]byte("addr"))

	_, ok := input.axk.GetAccountX(input.ctx, addr)
	require.False(t, ok)

	accx := input.axk.GetOrCreateAccountX(input.ctx, addr)
	require.Equal(t, addr, accx.Address)

	accx, ok = input.axk.GetAccountX(input.ctx, addr)
	require.True(t, ok)
	require.Equal(t, addr, accx.Address)
}
