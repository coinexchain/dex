package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountX_GetAllUnlockedCoinsAtTheTime(t *testing.T) {
	var acc = AccountX{ Address: []byte("123"), Activated: true, MemoRequired: false }
	coins := LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom:"cet", Amount:sdk.NewInt(20)}, UnlockTime: 1000},
		LockedCoin{Coin: sdk.Coin{Denom:"eth", Amount:sdk.NewInt(30)}, UnlockTime: 2000},
		LockedCoin{Coin: sdk.Coin{Denom:"eos", Amount:sdk.NewInt(40)}, UnlockTime: 3000},
	}
	acc.LockedCoins = coins
	res := acc.GetAllUnlockedCoinsAtTheTime(1000)
	require.Equal(t, res, LockedCoin{sdk.Coin{Denom:"cet", Amount:sdk.NewInt(20)}, 1000})
}
