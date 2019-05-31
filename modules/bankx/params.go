package bankx

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = "bankx"
)

var _ params.ParamSet = &Params{}

var (
	KeyActivationFee = []byte("ActivationFee")
)

type Params struct {
	ActivationFee int64 `json:"activation_fee"`
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyActivationFee, Value: &p.ActivationFee},
	}
}

func DefaultParams() Params {
	return Params{
		ActivationFee: 100000000,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
