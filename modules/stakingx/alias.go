package stakingx

import "github.com/coinexchain/dex/modules/stakingx/internal/types"

type (
	Params = types.Params
)

const (
	StoreKey                 = types.StoreKey
	ModuleName               = types.ModuleName
	QuerierRoute             = types.QuerierRoute
	DefaultParamspace        = types.DefaultParamspace
	DefaultMinSelfDelegation = types.DefaultMinSelfDelegation
)

var (
	DefaultParams                     = types.DefaultParams
	DefaultMinMandatoryCommissionRate = types.DefaultMinMandatoryCommissionRate
)
