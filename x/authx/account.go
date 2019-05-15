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

func NewAccountXWithAddress(ctx sdk.Context, addr sdk.AccAddress) AccountX{
	acc:=AccountX{
		Address:addr,
	}
	return acc
}
func (acc *AccountX) SetActivated(activated bool) {
	acc.Activated=activated
}