package authx

import (
	"github.com/coinexchain/dex/modules/authx/types"
)

const (
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	ModuleName   = types.ModuleName

	DefaultParamspace       = types.DefaultParamspace
	DefaultMinGasPriceLimit = types.DefaultMinGasPriceLimit
)

var (
	ErrInvalidMinGasPriceLimit = types.ErrInvalidMinGasPriceLimit
	ErrGasPriceTooLow          = types.ErrGasPriceTooLow
)

type (
	AccountX    = types.AccountX
	LockedCoin  = types.LockedCoin
	LockedCoins = types.LockedCoins
)
