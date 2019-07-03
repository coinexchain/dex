package authx

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestToString(t *testing.T) {
	lockedCoin := NewLockedCoin("cet", sdk.NewInt(100), 12345)
	require.Equal(t, "coin: 100cet, unlocked_time: 12345\n", lockedCoin.String())

	lockedCoins := LockedCoins{lockedCoin, lockedCoin}
	require.Equal(t, "coin: 100cet, unlocked_time: 12345\ncoin: 100cet, unlocked_time: 12345", lockedCoins.String())
}
