package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-----------------------------------------------------------------------------
// Locked Coin

type LockedCoin struct {
	Coin       sdk.Coin `json:"coin"`
	UnlockTime int64    `json:"unlock_time"`
}

func NewLockedCoin(denom string, amount sdk.Int, unlockTime int64) LockedCoin {
	return LockedCoin{
		Coin:       sdk.NewCoin(denom, amount),
		UnlockTime: unlockTime,
	}
}

func (coin LockedCoin) String() string {
	return fmt.Sprintf("coin: %s, unlocked_time: %d\n", coin.Coin, coin.UnlockTime)
}

//-----------------------------------------------------------------------------
// Locked Coins

type LockedCoins []LockedCoin

func (coins LockedCoins) String() string {
	out := ""
	if len(coins) == 0 {
		return out
	}

	for _, coin := range coins {
		out += coin.String()
	}
	return out[:len(out)-1]
}
