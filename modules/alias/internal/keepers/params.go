package keepers

import (
	"bytes"
	"fmt"

	"github.com/coinexchain/dex/modules/alias/internal/types"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultFeeForAliasLength2         = 2000000
	DefaultFeeForAliasLength3         = 300000
	DefaultFeeForAliasLength4         = 40000
	DefaultFeeForAliasLength5         = 5000
	DefaultFeeForAliasLength6         = 600
	DefaultFeeForAliasLength7OrHigher = 70
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

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
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

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
