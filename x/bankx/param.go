package bankx

import "github.com/cosmos/cosmos-sdk/x/params"

const (
	DefaultParamSpace = "bankx"
)

var ParamStoreKeyActivatedFee = []byte("ActivatedFee")

type Param struct {
	ActivatedFee int64
}

func DefaultParam() Param {
	return Param{
		ActivatedFee: 1,
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterType(ParamStoreKeyActivatedFee, &Param{})
}
