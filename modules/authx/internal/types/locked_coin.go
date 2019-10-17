package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-----------------------------------------------------------------------------
// Locked Coin

type LockedCoin struct {
	FromAddress sdk.AccAddress `json:"from_address,omitempty"`
	Supervisor  sdk.AccAddress `json:"supervisor,omitempty"`
	Coin        sdk.Coin       `json:"coin"`
	UnlockTime  int64          `json:"unlock_time"`
	Reward      int64          `json:"reward"`
}

func NewLockedCoin(fromAddress sdk.AccAddress, supervisor sdk.AccAddress, denom string, amount sdk.Int,
	unlockTime int64, reward int64) LockedCoin {
	return LockedCoin{
		FromAddress: fromAddress,
		Supervisor:  supervisor,
		Coin:        sdk.NewCoin(denom, amount),
		UnlockTime:  unlockTime,
		Reward:      reward,
	}
}

func (coin LockedCoin) String() string {
	return fmt.Sprintf("from: %s, supervisor: %s, coin: %s, unlocked_time: %d, reward: %d\n",
		coin.FromAddress.String(), coin.Supervisor.String(), coin.Coin, coin.UnlockTime, coin.Reward)
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
