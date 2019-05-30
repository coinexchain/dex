package authx

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default authx module parameter subspace
const DefaultParamspace = "authx"

// Default parameter values
const (
	DefaultMinSelfDelegation        = 1
	DefaultMinGasPrice       uint64 = 1e8 // 1 CET/Gas

	DefaultMsgSendGasCost uint64 = 1
)

// Parameter keys
var (
	KeyMinSelfDelegation = []byte("MinSelfDelegation")
	KeyMinGasPrice       = []byte("MinGasPrice")
	KeyMsgSendGasCost    = []byte("MsgSendGasCost")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the authx module.
type Params struct {
	MinSelfDelegation sdk.Int `json:"min_self_delegation"`
	MinGasPrice       uint64  `json:"min_gas_price"`
	MsgSendGasCost    uint64  `json:"msg_send_gas_cost"`
}

// ParamKeyTable for authx module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of authx module's parameters.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinSelfDelegation, Value: &p.MinSelfDelegation},
		{Key: KeyMinGasPrice, Value: &p.MinGasPrice},
		{Key: KeyMsgSendGasCost, Value: &p.MsgSendGasCost},
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
		MinSelfDelegation: sdk.NewInt(DefaultMinSelfDelegation),
		MinGasPrice:       DefaultMinGasPrice,
		MsgSendGasCost:    DefaultMsgSendGasCost,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	return string(msgCdc.MustMarshalBinaryLengthPrefixed(&p))
}
