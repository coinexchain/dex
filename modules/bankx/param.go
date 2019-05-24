package bankx

import "github.com/cosmos/cosmos-sdk/x/params"

const (
	DefaultParamspace = "bankx"
)

var ParamStoreKeyActivatedFee = []byte("ActivatedFee")

type Param struct {
	ActivatedFee int64 `json:"activated_fee"`
}

func DefaultParam() Param {
	return Param{
		ActivatedFee: 100000000,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterType(ParamStoreKeyActivatedFee, &Param{})
}
