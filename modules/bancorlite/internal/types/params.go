package types

import (
	"bytes"
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultCreateBancorFee = 1e10 // 100 * 10 ^8
	DefaultCancelBancorFee = 1e10 // 100 * 10 ^8
	TradeFeeRatePrecision  = 4
	DefaultTradeFeeRate    = 10
)

var (
	KeyCreateBancorFee = []byte("CreateBancorFee")
	KeyCancelBancorFee = []byte("CancelBancorFee")
	KeyTradeFeeRate    = []byte("TradeFeeRate")
)

type Params struct {
	CreateBancorFee int64 `json:"create_bancor_fee"`
	CancelBancorFee int64 `json:"cancel_bancor_fee"`
	TradeFeeRate    int64 `json:"trade_fee_rate"`
}

// ParamKeyTable for bancorlite module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		DefaultCreateBancorFee,
		DefaultCancelBancorFee,
		DefaultTradeFeeRate,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyCreateBancorFee, Value: &p.CreateBancorFee},
		{Key: KeyCancelBancorFee, Value: &p.CancelBancorFee},
		{Key: KeyTradeFeeRate, Value: &p.TradeFeeRate},
	}
}

func (p *Params) ValidateGenesis() error {
	if p.CreateBancorFee <= 0 || p.CancelBancorFee <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyCreateBancorFee, p.CreateBancorFee)
	}
	if p.TradeFeeRate < 0 || p.TradeFeeRate >= int64(math.Pow10(TradeFeeRatePrecision)) {
		return fmt.Errorf("TradeFeeRate is invalid")
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
	return fmt.Sprintf(`BancorLite Params:
  CreateBancorFee: %d
  CancelBancorFee: %d
  TradeFeeRate:    %d`,
		p.CreateBancorFee,
		p.CancelBancorFee,
		p.TradeFeeRate)
}
