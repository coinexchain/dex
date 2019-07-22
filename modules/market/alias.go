package market

import (
	"github.com/coinexchain/dex/modules/market/internal/types"

	"github.com/coinexchain/dex/modules/market/internal/keepers"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

const (
	IntegrationNetSubString = types.IntegrationNetSubString
	OrderIDPartsNum         = types.OrderIDPartsNum
	SymbolSeparator         = types.SymbolSeparator
	LimitOrder              = types.LimitOrder
	SELL                    = types.SELL
	GTE                     = types.GTE
)

var (
	NewBaseKeeper = keepers.NewKeeper
	DefaultParams = keepers.DefaultParams
)

type (
	Keeper               = keepers.Keeper
	Order                = types.Order
	MarketInfo           = types.MarketInfo
	MsgCreateOrder       = types.MsgCreateOrder
	MsgCreateTradingPair = types.MsgCreateTradingPair
)
