package authx

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/authx/internal/keepers"
	"github.com/coinexchain/dex/modules/authx/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	StoreKey        = types.StoreKey
	QuerierRoute    = types.QuerierRoute
	ModuleName      = types.ModuleName
	QueryAccountMix = types.QueryAccountMix

	CodeSpaceAuthX     = types.CodeSpaceAuthX
	CodeGasPriceTooLow = types.CodeGasPriceTooLow

	DefaultParamspace       = types.DefaultParamspace
	DefaultMinGasPriceLimit = types.DefaultMinGasPriceLimit
)

var (
	ErrInvalidMinGasPriceLimit = types.ErrInvalidMinGasPriceLimit
	ErrGasPriceTooLow          = types.ErrGasPriceTooLow
	NewLockedCoin              = types.NewLockedCoin
	NewSupervisedLockedCoin    = types.NewSupervisedLockedCoin
	NewParams                  = types.NewParams
	NewAccountX                = types.NewAccountX
	DefaultParams              = types.DefaultParams
	NewAccountXWithAddress     = types.NewAccountXWithAddress
	NewKeeper                  = keepers.NewKeeper
)

type (
	AccountX              = types.AccountX
	LockedCoin            = types.LockedCoin
	LockedCoins           = types.LockedCoins
	AccountXKeeper        = keepers.AccountXKeeper
	ExpectedAccountKeeper = keepers.ExpectedAccountKeeper
	ExpectedTokenKeeper   = keepers.ExpectedTokenKeeper
)
