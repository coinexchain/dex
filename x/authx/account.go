package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountX struct {
	Address              sdk.AccAddress `json:"address"`
	Activated            bool           `json:"activated"`
	TransferMemoRequired bool           `json:"transfer_memo_required"`
	LockedCoins          []LockedCoin   `json:"locked_coins"`
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

func (acc *AccountX) GetActivated() bool {
	return acc.Activated
}
func (acc *AccountX) SetActivated(activated bool) {
	acc.Activated = activated
}
func (acc *AccountX) SetTransferMemoRequired(transferMemoRequired bool) {
	acc.TransferMemoRequired = transferMemoRequired
}

func (acc *AccountX) GetTransferMemoRequired() bool {
	return acc.TransferMemoRequired
}
func (acc *AccountX) GetAddress() sdk.AccAddress {
	return acc.Address
}
func (acc *AccountX) GetLockedCoins() []LockedCoin {
	return acc.LockedCoins
}
