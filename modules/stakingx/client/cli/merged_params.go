package cli

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// staking.Params & stakingx.Params
type MergedParams struct {
	UnbondingTime              time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	MaxValidators              uint16        `json:"max_validators" yaml:"max_validators"`
	MaxEntries                 uint16        `json:"max_entries" yaml:"max_entries"`
	BondDenom                  string        `json:"bond_denom" yaml:"bond_denom"`
	MinSelfDelegation          sdk.Int       `json:"min_self_delegation" yaml:"min_self_delegation"`
	MinMandatoryCommissionRate sdk.Dec       `json:"min_mandatory_commission_rate" yaml:"min_mandatory_commission_rate"`
}

func (p MergedParams) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:                %s
  Max Validators:                %d
  Max Entries:                   %d
  Bonded Coin Denom:             %s
  Min Self Delegation:           %s
  Min Mandatory Commission Rate: %s`,
		p.UnbondingTime, p.MaxValidators, p.MaxEntries, p.BondDenom,
		p.MinSelfDelegation, p.MinMandatoryCommissionRate)
}
