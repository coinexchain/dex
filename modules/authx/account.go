package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type AccountX struct {
	Address      sdk.AccAddress `json:"address"`
	Activated    bool           `json:"activated"`
	MemoRequired bool           `json:"memo_required"` // if memo is required for receiving coins
	LockedCoins  []LockedCoin   `json:"locked_coins"`
}

func (acc *AccountX) GetAddress() sdk.AccAddress {
	return acc.Address
}

func (acc *AccountX) SetAddress(address sdk.AccAddress) {
	acc.Address = address
}

func (acc *AccountX) SetMemoRequired(b bool) {
	acc.MemoRequired = b
}

func (acc *AccountX) IsMemoRequired() bool {
	return acc.MemoRequired
}

func (acc *AccountX) SetActivated(b bool) {
	acc.Activated = b
}

func (acc AccountX) IsActivated() bool {
	return acc.Activated
}

func (acc *AccountX) AddLockedCoins(coins LockedCoins) {
	acc.LockedCoins = append(acc.LockedCoins, coins...)
}

func (acc *AccountX) GetAllLockedCoins() LockedCoins {
	return acc.LockedCoins
}

func (acc *AccountX) GetLockedCoinsByDemon(demon string) LockedCoins {
	var coins LockedCoins
	for _, c := range acc.LockedCoins {
		if c.Coin.Denom == demon {
			coins = append(coins, c)
		}
	}
	return coins
}

func (acc *AccountX) GetUnlockedCoinsAtTheTime(demon string, time int64) LockedCoins {
	var coins LockedCoins
	for _, c := range acc.GetLockedCoinsByDemon(demon) {
		if c.UnlockTime <= time {
			coins = append(coins, c)
		}
	}
	return coins
}

func (acc *AccountX) GetAllUnlockedCoinsAtTheTime(time int64) LockedCoins {
	var coins LockedCoins
	for _, c := range acc.GetAllLockedCoins() {
		if c.UnlockTime <= time {
			coins = append(coins, c)
		}
	}
	return coins
}

func (acc *AccountX) TransferUnlockedCoins(time int64, ctx sdk.Context, kx AccountXKeeper, keeper auth.AccountKeeper) {
	account := keeper.GetAccount(ctx, acc.Address)
	oldCoins := account.GetCoins()
	var coins sdk.Coins
	var temp LockedCoins
	for _, c := range acc.LockedCoins {
		if c.UnlockTime <= time {
			coins = append(coins, c.Coin)
		} else {
			temp = append(temp, c)
		}
	}
	newCoins := oldCoins.Add(coins)
	account.SetCoins(newCoins)
	keeper.SetAccount(ctx, account)
	acc.LockedCoins = temp
	kx.SetAccountX(ctx, *acc)
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
