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
	KeyMinSelfDelegation          = []byte("MinSelfDelegation")
	KeyMinMandatoryCommissionRate = []byte("MinMandatoryCommissionRate")

	DefaultMinMandatoryCommissionRate = sdk.NewDecWithPrec(5, 2)
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the stakingx module.
type Params struct {
	MinSelfDelegation          sdk.Int `json:"min_self_delegation"`
	MinMandatoryCommissionRate sdk.Dec `json:"min_mandatory_commission_rate"`
}

// ParamKeyTable for stakingx module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MinSelfDelegation:          sdk.NewInt(DefaultMinSelfDelegation),
		MinMandatoryCommissionRate: DefaultMinMandatoryCommissionRate,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of stakingx module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinSelfDelegation, Value: &p.MinSelfDelegation},
		{Key: KeyMinMandatoryCommissionRate, Value: &p.MinMandatoryCommissionRate},
	}
}
