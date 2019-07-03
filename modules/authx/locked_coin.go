package authx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LockedCoins []LockedCoin

type LockedCoin struct {
	Coin       sdk.Coin `json:"coin"`
	UnlockTime int64    `json:"unlock_time"`
}

func (coins LockedCoins) String() string {

	if len(coins) == 0 {
		return ""
	}
	out := ""
	for i, coin := range coins {
		if i == 0 {
			out += fmt.Sprintf("%v", coin.String())
		} else {
			out += fmt.Sprintf("                 %v", coin.String())
			//out += fmt.Sprintf("%17v", coin.String())
		}
	}
	return out[:len(out)-1]
}

func (coin LockedCoin) String() string {
	return fmt.Sprintf("coin: %s, unlocked_time: %d\n", coin.Coin.String(), coin.UnlockTime)
}
