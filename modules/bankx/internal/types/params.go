package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = &Params{}

var (
	KeyActivationFee = []byte("ActivationFee")
	KeyLockCoinsFee  = []byte("LockCoinsFee")
)

type Params struct {
	ActivationFee int64 `json:"activation_fee"`
	LockCoinsFee  int64 `json:"lock_coins_fee"`
}

func DefaultParams() Params {
	return Params{
		ActivationFee: 100000000,
		LockCoinsFee:  1e10,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyActivationFee, Value: &p.ActivationFee},
		{Key: KeyLockCoinsFee, Value: &p.LockCoinsFee},
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`Alias Params:
  ActivationFee: %d
  LockCoinsFee:  %d`,
		p.ActivationFee,
		p.LockCoinsFee)
}
