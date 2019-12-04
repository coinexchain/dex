package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default asset module parameter subspace
const (
	MaxTokenAmount = 5e76 // 57896044618658097711785492504343953926634992332820282019728792003956564819967

	DefaultIssue2CharTokenFee = 100000e8 // 100000 * 10^8
	DefaultIssue3CharTokenFee = 10000e8  //  10000 * 10^8
	DefaultIssue4CharTokenFee = 5000e8   //   5000 * 10^8
	DefaultIssue5CharTokenFee = 2000e8   //   2000 * 10^8
	DefaultIssue6CharTokenFee = 1000e8   //   1000 * 10^8
	DefaultIssueLongTokenFee  = 500e8    //    500 * 10^8
)

// Parameter keys
var (
	KeyIssueTokenFee      = []byte("IssueTokenFee")
	KeyIssueRareTokenFee  = []byte("IssueRareTokenFee")
	KeyIssue3CharTokenFee = []byte("Issue3CharTokenFee")
	KeyIssue4CharTokenFee = []byte("Issue4CharTokenFee")
	KeyIssue5CharTokenFee = []byte("Issue5CharTokenFee")
	KeyIssue6CharTokenFee = []byte("Issue6CharTokenFee")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the parameters for the asset module.
type Params struct {
	// FeeParams define the rules according to which fee are charged.
	IssueTokenFee      int64 `json:"issue_token_fee" yaml:"issue_token_fee"`             // 7+ char
	IssueRareTokenFee  int64 `json:"issue_rare_token_fee" yaml:"issue_rare_token_fee"`   // 2 char
	Issue3CharTokenFee int64 `json:"issue_3char_token_fee" yaml:"issue_3char_token_fee"` // 3 char
	Issue4CharTokenFee int64 `json:"issue_4char_token_fee" yaml:"issue_4char_token_fee"` // 4 char
	Issue5CharTokenFee int64 `json:"issue_5char_token_fee" yaml:"issue_5char_token_fee"` // 5 char
	Issue6CharTokenFee int64 `json:"issue_6char_token_fee" yaml:"issue_6char_token_fee"` // 6 char
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		IssueTokenFee:      DefaultIssueLongTokenFee,
		IssueRareTokenFee:  DefaultIssue2CharTokenFee,
		Issue3CharTokenFee: DefaultIssue3CharTokenFee,
		Issue4CharTokenFee: DefaultIssue4CharTokenFee,
		Issue5CharTokenFee: DefaultIssue5CharTokenFee,
		Issue6CharTokenFee: DefaultIssue6CharTokenFee,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyIssueTokenFee, Value: &p.IssueTokenFee},
		{Key: KeyIssueRareTokenFee, Value: &p.IssueRareTokenFee},
		{Key: KeyIssue3CharTokenFee, Value: &p.Issue3CharTokenFee},
		{Key: KeyIssue4CharTokenFee, Value: &p.Issue4CharTokenFee},
		{Key: KeyIssue5CharTokenFee, Value: &p.Issue5CharTokenFee},
		{Key: KeyIssue6CharTokenFee, Value: &p.Issue6CharTokenFee},
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

func (p Params) GetIssueTokenFee(symbol string) int64 {
	switch len(symbol) {
	case 2:
		return p.IssueRareTokenFee
	case 3:
		return p.Issue3CharTokenFee
	case 4:
		return p.Issue4CharTokenFee
	case 5:
		return p.Issue5CharTokenFee
	case 6:
		return p.Issue6CharTokenFee
	default:
		return p.IssueTokenFee
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func (p Params) String() string {
	return fmt.Sprintf(`Asset Params:
  IssueTokenFee:      %d
  IssueRareTokenFee:  %d
  Issue3CharTokenFee: %d
  Issue4CharTokenFee: %d
  Issue5CharTokenFee: %d
  Issue6CharTokenFee: %d`,
		p.IssueTokenFee,
		p.IssueRareTokenFee,
		p.Issue3CharTokenFee,
		p.Issue4CharTokenFee,
		p.Issue5CharTokenFee,
		p.Issue6CharTokenFee,
	)
}
