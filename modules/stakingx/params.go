package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter values
const (
	DefaultMinSelfDelegation = 1000000e8
)

// Parameter keys
var (
	KeyMinSelfDelegation    = []byte("MinSelfDelegation")
	KeyNonBondableAddresses = []byte("NonBondableAddresses")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the stakingx module.
type Params struct {
	MinSelfDelegation    sdk.Int          `json:"min_self_delegation"`
	NonBondableAddresses []sdk.AccAddress `json:"non_bondable_addresses"`
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
		{Key: KeyNonBondableAddresses, Value: &p.NonBondableAddresses},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MinSelfDelegation:    sdk.NewInt(DefaultMinSelfDelegation),
		NonBondableAddresses: make([]sdk.AccAddress, 0),
	}
}
