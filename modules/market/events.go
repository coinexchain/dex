package market

// Market module event types
var (
	EventTypeMarket = "market"

	AttributeKeyTradingPair      = "trading-pair"
	AttributeKeyOrder            = "order"
	AttributeKeyStock            = "stock"
	AttributeKeyMoney            = "money"
	AttributeKeyPricePrecision   = "price-precision"
	AttributeKeyLastExecutePrice = "last-execute-price"

	AttributeKeySequence    = "sequence"
	AttributeKeyOrderType   = "order-type"
	AttributeKeyPrice       = "price"
	AttributeKeyQuantity    = "quantity"
	AttributeKeySide        = "side"
	AttributeKeyHeight      = "height"
	AttributeKeyTimeInForce = "time-in-force"

	AttributeKeyDelOrderReason = "del-order-reason"
	AttributeKeyDelOrderHeight = "del-order-height"
	AttributeKeyUsedCommission = "used-commission"
	AttributeKeyLeftStock      = "left-stock"
	AttributeKeyDealStock      = "deal-stock"
	AttributeKeyRemainAmount   = "remain-amount"
	AttributeKeyDealMoney      = "deal-money"

	AttributeKeyEffectiveTime = "effective-time"

	AttributeKeyOldPricePrecision = "old-price-precision"
	AttributeKeyNewPricePrecision = "new-price-precision"
)
