package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountX struct {
	Address      sdk.AccAddress `json:"address"`
	Activated    bool           `json:"activated"`
	MemoRequired bool           `json:"memo_required"` // if memo is required for receiving coins
	LockedCoins  []LockedCoin   `json:"locked_coins"`
}

type LockedCoin struct {
	Coin       sdk.Coin `json:"coin"`
	UnlockTime int64    `json:"unlock_time"`
}

func NewAccountXWithAddress(addr sdk.AccAddress) AccountX {
	acc := AccountX{
		Address: addr,
	}
	return acc
}
