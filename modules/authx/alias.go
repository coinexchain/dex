package authx

import (
	"github.com/coinexchain/dex/modules/authx/types"
)

const (
	StoreKey      = types.StoreKey
	QuerierRoute  = types.QuerierRoute
	ModuleName    = types.ModuleName
	QueryAccountX = types.QueryAccountX

	DefaultParamspace       = types.DefaultParamspace
	DefaultMinGasPriceLimit = types.DefaultMinGasPriceLimit
)

var (
	ErrInvalidMinGasPriceLimit = types.ErrInvalidMinGasPriceLimit
	ErrGasPriceTooLow          = types.ErrGasPriceTooLow
	NewLockedCoin              = types.NewLockedCoin
	NewParams                  = types.NewParams
)

type (
	AccountX    = types.AccountX
	LockedCoin  = types.LockedCoin
	LockedCoins = types.LockedCoins
)
