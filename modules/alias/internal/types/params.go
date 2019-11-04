package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultFeeForAliasLength2         = 10000e8
	DefaultFeeForAliasLength3         = 5000e8
	DefaultFeeForAliasLength4         = 2000e8
	DefaultFeeForAliasLength5         = 1000e8
	DefaultFeeForAliasLength6         = 100e8
	DefaultFeeForAliasLength7OrHigher = 10e8
	DefaultMaxAliasCount              = 5
)

var (
	KeyFeeForAliasLength2         = []byte("FeeForAliasLength2")
	KeyFeeForAliasLength3         = []byte("FeeForAliasLength3")
	KeyFeeForAliasLength4         = []byte("FeeForAliasLength4")
	KeyFeeForAliasLength5         = []byte("FeeForAliasLength5")
	KeyFeeForAliasLength6         = []byte("FeeForAliasLength6")
	KeyFeeForAliasLength7OrHigher = []byte("FeeForAliasLength7OrHigher")
	KeyMaxAliasCount              = []byte("MaxAliasCount")
)

type Params struct {
	FeeForAliasLength2         int64 `json:"fee_for_alias_length_2"`
	FeeForAliasLength3         int64 `json:"fee_for_alias_length_3"`
	FeeForAliasLength4         int64 `json:"fee_for_alias_length_4"`
	FeeForAliasLength5         int64 `json:"fee_for_alias_length_5"`
	FeeForAliasLength6         int64 `json:"fee_for_alias_length_6"`
	FeeForAliasLength7OrHigher int64 `json:"fee_for_alias_length_7_or_higher"`
	MaxAliasCount              int   `json:"max_alias_count"`
}

// ParamKeyTable for alias module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		DefaultFeeForAliasLength2,
		DefaultFeeForAliasLength3,
		DefaultFeeForAliasLength4,
		DefaultFeeForAliasLength5,
		DefaultFeeForAliasLength6,
		DefaultFeeForAliasLength7OrHigher,
		DefaultMaxAliasCount,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyFeeForAliasLength2, Value: &p.FeeForAliasLength2},
		{Key: KeyFeeForAliasLength3, Value: &p.FeeForAliasLength3},
		{Key: KeyFeeForAliasLength4, Value: &p.FeeForAliasLength4},
		{Key: KeyFeeForAliasLength5, Value: &p.FeeForAliasLength5},
		{Key: KeyFeeForAliasLength6, Value: &p.FeeForAliasLength6},
		{Key: KeyFeeForAliasLength7OrHigher, Value: &p.FeeForAliasLength7OrHigher},
		{Key: KeyMaxAliasCount, Value: &p.MaxAliasCount},
	}
}

func (p *Params) ValidateGenesis() error {
	if p.FeeForAliasLength2 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength2, p.FeeForAliasLength2)
	}
	if p.FeeForAliasLength3 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength3, p.FeeForAliasLength3)
	}
	if p.FeeForAliasLength3 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength3, p.FeeForAliasLength3)
	}
	if p.FeeForAliasLength4 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength4, p.FeeForAliasLength4)
	}
	if p.FeeForAliasLength5 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength5, p.FeeForAliasLength5)
	}
	if p.FeeForAliasLength6 <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength6, p.FeeForAliasLength6)
	}
	if p.FeeForAliasLength7OrHigher <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyFeeForAliasLength7OrHigher, p.FeeForAliasLength7OrHigher)
	}
	if p.MaxAliasCount <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyMaxAliasCount, p.MaxAliasCount)
	}
	return nil
}

func (p *Params) GetFeeForAlias(alias string) int64 {
	if n := len(alias); n == 2 {
		return p.FeeForAliasLength2
	} else if n == 3 {
		return p.FeeForAliasLength3
	} else if n == 4 {
		return p.FeeForAliasLength4
	} else if n == 5 {
		return p.FeeForAliasLength5
	} else if n == 6 {
		return p.FeeForAliasLength6
	} else {
		return p.FeeForAliasLength7OrHigher
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func (p Params) String() string {
	return fmt.Sprintf(`Alias Params:
  FeeForAliasLength2:         %d
  FeeForAliasLength3:         %d
  FeeForAliasLength4:         %d
  FeeForAliasLength5:         %d
  FeeForAliasLength6:         %d
  FeeForAliasLength7OrHigher: %d
  MaxAliasCount:              %d`,
		p.FeeForAliasLength2,
		p.FeeForAliasLength3,
		p.FeeForAliasLength4,
		p.FeeForAliasLength5,
		p.FeeForAliasLength6,
		p.FeeForAliasLength7OrHigher,
		p.MaxAliasCount)
}
