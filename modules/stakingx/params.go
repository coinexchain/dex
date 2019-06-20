package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default stakingx module parameter subspace
const DefaultParamspace = "stakingx"

// Default parameter values
const (
	DefaultMinSelfDelegation = 10000e8
)

// Parameter keys
var (
	KeyMinSelfDelegation = []byte("MinSelfDelegation")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the stakingx module.
type Params struct {
	MinSelfDelegation sdk.Int `json:"min_self_delegation"`
}

// ParamKeyTable for stakingx module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of stakingx module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinSelfDelegation, Value: &p.MinSelfDelegation},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MinSelfDelegation: sdk.NewInt(DefaultMinSelfDelegation),
	}
}
