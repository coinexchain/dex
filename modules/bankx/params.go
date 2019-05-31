package bankx

import "github.com/cosmos/cosmos-sdk/x/params"

const (
	DefaultParamspace = "bankx"
)

var ParamStoreKeyActivationFee = []byte("ActivationFee")

type Params struct {
	ActivationFee int64 `json:"activation_fee"`
}

func DefaultParams() Params {
	return Params{
		ActivationFee: 100000000,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterType(ParamStoreKeyActivationFee, &Params{})
}
