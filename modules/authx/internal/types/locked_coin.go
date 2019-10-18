package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-----------------------------------------------------------------------------
// Locked Coin

type LockedCoin struct {
	Coin        sdk.Coin       `json:"coin"`
	UnlockTime  int64          `json:"unlock_time"`
	FromAddress sdk.AccAddress `json:"from_address,omitempty"`
	Supervisor  sdk.AccAddress `json:"supervisor,omitempty"`
	Reward      int64          `json:"reward,omitempty"`
}

func NewLockedCoin(denom string, amount sdk.Int, unlockTime int64) LockedCoin {
	return LockedCoin{
		Coin:       sdk.NewCoin(denom, amount),
		UnlockTime: unlockTime,
	}
}

func NewSupervisedLockedCoin(denom string, amount sdk.Int, unlockTime int64, fromAddress sdk.AccAddress, supervisor sdk.AccAddress, reward int64) LockedCoin {
	return LockedCoin{
		Coin:        sdk.NewCoin(denom, amount),
		UnlockTime:  unlockTime,
		FromAddress: fromAddress,
		Supervisor:  supervisor,
		Reward:      reward,
	}
}

func (coin LockedCoin) String() string {
	str := fmt.Sprintf("coin: %s, unlocked_time: %d", coin.Coin, coin.UnlockTime)
	if coin.FromAddress != nil {
		str += fmt.Sprintf(", from: %s", coin.FromAddress.String())
	}
	if coin.Supervisor != nil {
		str += fmt.Sprintf(", supervisor: %s, reward: %d", coin.Supervisor.String(), coin.Reward)
	}
	str += "\n"
	return str
}

func (coin LockedCoin) IsEqual(other LockedCoin) bool {
	return coin.Coin.IsEqual(other.Coin) &&
		coin.UnlockTime == other.UnlockTime &&
		coin.Reward == other.Reward &&
		bytes.Equal(coin.FromAddress, other.FromAddress) &&
		bytes.Equal(coin.Supervisor, other.Supervisor)
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
