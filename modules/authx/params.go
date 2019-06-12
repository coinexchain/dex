package authx

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default authx module parameter subspace
const DefaultParamspace = "authx"

// Default parameter values
const (
	// DefaultMinGasPrice of the network
	// Make token transfer/send tx to costs around 0.01CET
	// activated account send to self,                  costs 38883 gas
	// activated account send to non-activated account, costs 48951 gas
	// activated account send to other activated addr,  costs 33903 gas
	// consider it takes 50000 to do transfer/send tx
	// so, min_gas_price = 100000000sato.CET * 0.01 / 50000 = 20 sato.CET
	DefaultMinGasPriceLimit uint64 = 20
)

// Parameter keys
var (
	KeyMinGasPriceLimit = []byte("MinGasPriceLimit")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the authx module.
type Params struct {
	MinGasPriceLimit uint64 `json:"min_gas_price_limit"`
}

// ParamKeyTable for authx module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of authx module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinGasPriceLimit, Value: &p.MinGasPriceLimit},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MinGasPriceLimit: DefaultMinGasPriceLimit,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	return string(msgCdc.MustMarshalBinaryLengthPrefixed(&p))
}
