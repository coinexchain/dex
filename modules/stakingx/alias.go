package stakingx

import (
	"github.com/coinexchain/dex/modules/stakingx/internal/keepers"
	"github.com/coinexchain/dex/modules/stakingx/internal/types"
)

type (
	Params = types.Params
)

const (
	StoreKey                            = types.StoreKey
	ModuleName                          = types.ModuleName
	QuerierRoute                        = types.QuerierRoute
	DefaultParamspace                   = types.DefaultParamspace
	DefaultMinSelfDelegation            = types.DefaultMinSelfDelegation
	CodeMinSelfDelegationBelowRequired  = types.CodeMinSelfDelegationBelowRequired
	CodeBelowMinMandatoryCommissionRate = types.CodeBelowMinMandatoryCommissionRate
)

type (
	Keeper = keepers.Keeper
)

var (
	DefaultParams                          = types.DefaultParams
	DefaultMinMandatoryCommissionRate      = types.DefaultMinMandatoryCommissionRate
	ErrInvalidMinSelfDelegation            = types.ErrInvalidMinSelfDelegation
	ErrMinSelfDelegationBelowRequired      = types.ErrMinSelfDelegationBelowRequired
	ErrRateBelowMinMandatoryCommissionRate = types.ErrRateBelowMinMandatoryCommissionRate
	NewKeeper                              = keepers.NewKeeper
)
