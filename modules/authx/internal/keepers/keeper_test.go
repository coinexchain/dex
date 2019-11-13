package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/authx/internal/keepers"
	"github.com/coinexchain/dex/modules/authx/internal/types"
)

var addr = sdk.AccAddress([]byte("some-address"))

func TestGetSetParams(t *testing.T) {
	codec.RunInitFuncList()
	input := setupTestInput()
	params := types.DefaultParams()
	input.axk.SetParams(input.ctx, params)
	params2 := input.axk.GetParams(input.ctx)
	require.True(t, params.Equal(params2))
}

func TestAccountXGetSet(t *testing.T) {
	codec.RunInitFuncList()
	input := setupTestInput()

	_, ok := input.axk.GetAccountX(input.ctx, addr)
	require.False(t, ok)

	//create account
	acc := types.NewAccountXWithAddress(addr)
	require.Equal(t, addr, acc.Address)

	input.axk.SetAccountX(input.ctx, acc)

	acc, ok = input.axk.GetAccountX(input.ctx, addr)
	require.True(t, ok)

	acc.MemoRequired = false
	input.axk.SetAccountX(input.ctx, acc)
	acc, _ = input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t, false, acc.MemoRequired)

	lockedCoins := acc.LockedCoins
	require.Nil(t, lockedCoins)
}

func TestAddressStoreKey(t *testing.T) {
	codec.RunInitFuncList()
	addrStoreKey := keepers.AddressStoreKey(addr)
	expectedOutput := []byte{0x1, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73}
	require.Equal(t, expectedOutput, addrStoreKey)
}

func TestGetOrCreateAccountX(t *testing.T) {
	codec.RunInitFuncList()
	input := setupTestInput()

	_, ok := input.axk.GetAccountX(input.ctx, addr)
	require.False(t, ok)

	accx := input.axk.GetOrCreateAccountX(input.ctx, addr)
	require.Equal(t, addr, accx.Address)

	accx, ok = input.axk.GetAccountX(input.ctx, addr)
	require.True(t, ok)
	require.Equal(t, addr, accx.Address)
}

func TestIteratorAccounts(t *testing.T) {
	codec.RunInitFuncList()
	input := setupTestInput()

	input.axk.GetOrCreateAccountX(input.ctx, sdk.AccAddress([]byte("addr0")))
	input.axk.GetOrCreateAccountX(input.ctx, sdk.AccAddress([]byte("addr1")))
	input.axk.GetOrCreateAccountX(input.ctx, sdk.AccAddress([]byte("addr2")))
	input.axk.GetOrCreateAccountX(input.ctx, sdk.AccAddress([]byte("addr3")))

	var accxs []types.AccountX
	input.axk.IterateAccounts(input.ctx, func(accx types.AccountX) bool {
		accxs = append(accxs, accx)
		return false
	})

	require.Equal(t, 4, len(accxs))
}
