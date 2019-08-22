package types

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type AccountX struct {
	Address      sdk.AccAddress `json:"address"`
	MemoRequired bool           `json:"memo_required"` // if memo is required for receiving coins
	LockedCoins  LockedCoins    `json:"locked_coins"`
	FrozenCoins  sdk.Coins      `json:"frozen_coins"`
}

type AccountXs []AccountX

func NewAccountX(address sdk.AccAddress, memoRequired bool, lockedCoins LockedCoins, frozenCoins sdk.Coins) AccountX {
	return AccountX{
		Address:      address,
		MemoRequired: memoRequired,
		LockedCoins:  lockedCoins,
		FrozenCoins:  frozenCoins,
	}
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

func (acc *AccountX) AddLockedCoins(coins LockedCoins) {
	acc.LockedCoins = append(acc.LockedCoins, coins...)
}

func (acc *AccountX) GetAllLockedCoins() LockedCoins {
	return acc.LockedCoins
}

func (acc *AccountX) GetAllCoins() sdk.Coins {
	lockedCoins := acc.LockedCoins
	coins := sdk.Coins{}
	for _, lockedCoin := range lockedCoins {
		coins = coins.Add(sdk.Coins{lockedCoin.Coin})
	}
	coins = coins.Add(acc.FrozenCoins)
	return coins
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

func (acc AccountX) String() string {
	return fmt.Sprintf(`
  LockedCoins:   %s
  FrozenCoins:   %s
  MemoRequired:  %t`,
		acc.LockedCoins, acc.FrozenCoins, acc.MemoRequired,
	)
}

func NewAccountXWithAddress(addr sdk.AccAddress) AccountX {
	acc := AccountX{
		Address: addr,
	}
	return acc
}

type AccountAll struct {
	Account  auth.Account `json:"account"`
	AccountX AccountX     `json:"account_x"`
}

func (accAll AccountAll) String() string {
	return accAll.Account.String() + accAll.AccountX.String()
}

type AccountMix struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	LockedCoins   LockedCoins    `json:"locked_coins"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
	MemoRequired  bool           `json:"memo_required"` // if memo is required for receiving coins
}

func NewAccountMix(acc auth.Account, x AccountX) AccountMix {
	return AccountMix{
		acc.GetAddress(),
		acc.GetCoins(),
		x.GetAllLockedCoins(),
		x.FrozenCoins,
		acc.GetPubKey(),
		acc.GetAccountNumber(),
		acc.GetSequence(),
		x.IsMemoRequired()}
}
