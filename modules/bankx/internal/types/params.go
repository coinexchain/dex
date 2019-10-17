package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

var (
	KeyActivationFee      = []byte("ActivationFee")
	KeyLockCoinsFreeTime  = []byte("LockCoinsFreeTime")
	KeyLockCoinsFeePerDay = []byte("LockCoinsFeePerDay")
)

type Params struct {
	ActivationFee      int64 `json:"activation_fee"`
	LockCoinsFreeTime  int64 `json:"lock_coins_free_time"`
	LockCoinsFeePerDay int64 `json:"lock_coins_fee_per_day"`
}

func NewParams(activation int64, freeTime, lock int64) Params {
	return Params{
		ActivationFee:      activation,
		LockCoinsFreeTime:  freeTime,
		LockCoinsFeePerDay: lock,
	}
}
func DefaultParams() Params {
	return Params{
		ActivationFee:      100000000,
		LockCoinsFreeTime:  604800000000000,
		LockCoinsFeePerDay: 1e6,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyActivationFee, Value: &p.ActivationFee},
		{Key: KeyLockCoinsFreeTime, Value: &p.LockCoinsFreeTime},
		{Key: KeyLockCoinsFeePerDay, Value: &p.LockCoinsFeePerDay},
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`BankX Params:
  ActivationFee:      %d
  LockCoinsFreeTime:  %d
  LockCoinsFeePerDay: %d`,
		p.ActivationFee,
		p.LockCoinsFreeTime,
		p.LockCoinsFeePerDay)
}
