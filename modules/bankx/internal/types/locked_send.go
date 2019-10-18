package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LockedSendMsg struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
	UnlockTime  int64          `json:"unlock_time"`
	Supervisor  sdk.AccAddress `json:"supervisor,omitempty"`
	Reward      int64          `json:"reward,omitempty"`
}

func NewLockedSendMsg(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amount sdk.Coins, unlockTime int64) LockedSendMsg {
	return LockedSendMsg{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
		UnlockTime:  unlockTime,
	}
}

func NewSupervisedSendMsg(fromAddr, toAddr, supervisor sdk.AccAddress, amount sdk.Coin, unlockTime int64, reward int64) LockedSendMsg {
	return LockedSendMsg{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      sdk.NewCoins(amount),
		UnlockTime:  unlockTime,
		Supervisor:  supervisor,
		Reward:      reward,
	}
}
