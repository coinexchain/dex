package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestToString(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	fromAddr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	supervisor, _ := sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	lockedCoin := NewLockedCoin(fromAddr, supervisor, "cet", sdk.NewInt(100), 12345, 1)
	require.Equal(t,
		"from: coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a, supervisor: coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd, coin: 100cet, unlocked_time: 12345, reward: 1\n",
		lockedCoin.String())

	lockedCoin2 := NewLockedCoin(nil, nil, "cet", sdk.NewInt(100), 12345, 0)
	require.Equal(t,
		"from: , supervisor: , coin: 100cet, unlocked_time: 12345, reward: 0\n",
		lockedCoin2.String())

	lockedCoins := LockedCoins{lockedCoin, lockedCoin2}
	require.Equal(t,
		"from: coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a, supervisor: coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd, coin: 100cet, unlocked_time: 12345, reward: 1\n"+
			"from: , supervisor: , coin: 100cet, unlocked_time: 12345, reward: 0",
		lockedCoins.String())
}
