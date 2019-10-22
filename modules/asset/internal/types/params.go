package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default asset module parameter subspace
const (
	MaxTokenAmount   = 5e76 // 57896044618658097711785492504343953926634992332820282019728792003956564819967
	RareSymbolLength = 2

	DefaultIssueTokenFee     = 1e12 // 10000 * 10 ^8
	DefaultIssueRareTokenFee = 1e13 // 100000 * 10 ^8
)

// Parameter keys
var (
	KeyIssueTokenFee     = []byte("IssueTokenFee")
	KeyIssueRareTokenFee = []byte("IssueRareTokenFee")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the parameters for the asset module.
type Params struct {
	// FeeParams define the rules according to which fee are charged.
	IssueTokenFee     int64 `json:"issue_token_fee" yaml:"issue_token_fee"`
	IssueRareTokenFee int64 `json:"issue_rare_token_fee" yaml:"issue_rare_token_fee"`
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		IssueTokenFee:     DefaultIssueTokenFee,
		IssueRareTokenFee: DefaultIssueRareTokenFee,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyIssueTokenFee, Value: &p.IssueTokenFee},
		{Key: KeyIssueRareTokenFee, Value: &p.IssueRareTokenFee},
	}
}

func (p *Params) ValidateGenesis() error {
	for _, pair := range p.ParamSetPairs() {
		fee := *(pair.Value.(*int64))
		if fee <= 0 {
			return fmt.Errorf("%s is invalid: %d", pair.Key, fee)
		}
	}
	return nil
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func (p Params) String() string {
	return fmt.Sprintf(`Asset Params:
  IssueTokenFee:     %d
  IssueRareTokenFee: %d`,
		p.IssueTokenFee,
		p.IssueRareTokenFee)
}
