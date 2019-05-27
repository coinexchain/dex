package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx"
)

// GenesisAccount defines an account initialized at genesis.
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting"`  // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free"`    // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time"`        // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time"`          // vesting end time (UNIX Epoch time)

	//accountx fields
	Activated    bool               `json:"activated"`
	MemoRequired bool               `json:"memo_required"` // if memo is required for receiving coins
	LockedCoins  []authx.LockedCoin `json:"locked_coins"`
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() auth.Account {
	bacc := &auth.BaseAccount{
		Address:       ga.Address,
		Coins:         ga.Coins.Sort(),
		AccountNumber: ga.AccountNumber,
		Sequence:      ga.Sequence,
	}

	if !ga.OriginalVesting.IsZero() {
		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:      bacc,
			OriginalVesting:  ga.OriginalVesting,
			DelegatedFree:    ga.DelegatedFree,
			DelegatedVesting: ga.DelegatedVesting,
			EndTime:          ga.EndTime,
		}

		if ga.StartTime != 0 && ga.EndTime != 0 {
			return &auth.ContinuousVestingAccount{
				BaseVestingAccount: baseVestingAcc,
				StartTime:          ga.StartTime,
			}
		} else if ga.EndTime != 0 {
			return &auth.DelayedVestingAccount{
				BaseVestingAccount: baseVestingAcc,
			}
		} else {
			panic(fmt.Sprintf("invalid genesis vesting account: %+v", ga))
		}
	}

	return bacc
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		Activated:     true,
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		gacc.OriginalVesting = vacc.GetOriginalVesting()
		gacc.DelegatedFree = vacc.GetDelegatedFree()
		gacc.DelegatedVesting = vacc.GetDelegatedVesting()
		gacc.StartTime = vacc.GetStartTime()
		gacc.EndTime = vacc.GetEndTime()
	}

	gacc.Activated = true
	return gacc
}
