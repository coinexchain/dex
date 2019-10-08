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
	GTE                     = types.GTE
	BID                     = types.BID
	ASK                     = types.ASK
	BUY                     = types.BUY
	SELL                    = types.SELL
)

var (
	NewBaseKeeper       = keepers.NewKeeper
	DefaultParams       = types.DefaultParams
	DecToBigEndianBytes = types.DecToBigEndianBytes
	ValidateOrderID     = types.ValidateOrderID
	IsValidTradingPair  = types.IsValidTradingPair
	ModuleCdc           = types.ModuleCdc
	GetSymbol           = types.GetSymbol
	SplitSymbol         = types.SplitSymbol
)

type (
	Keeper                  = keepers.Keeper
	Order                   = types.Order
	MarketInfo              = types.MarketInfo
	MsgCreateOrder          = types.MsgCreateOrder
	MsgCreateTradingPair    = types.MsgCreateTradingPair
	MsgCancelOrder          = types.MsgCancelOrder
	MsgCancelTradingPair    = types.MsgCancelTradingPair
	MsgModifyPricePrecision = types.MsgModifyPricePrecision
	CreateOrderInfo         = types.CreateOrderInfo
	FillOrderInfo           = types.FillOrderInfo
	CancelOrderInfo         = types.CancelOrderInfo
)
