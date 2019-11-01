package market

// Market module event types
var (
	AttributeValueCategory = ModuleName

	EventTypeKeyCreateTradingPair    = "create_market"
	EventTypeKeyCreateOrder          = "create_order"
	EventTypeKeyCancelOrder          = "cancel_order"
	EventTypeKeyCancelTradingPair    = "cancel_market"
	EventTypeKeyModifyPricePrecision = "modify_price_precision"

	AttributeKeyTradingPair      = "trading_pair"
	AttributeKeyOrder            = "order"
	AttributeKeyStock            = "stock"
	AttributeKeyMoney            = "money"
	AttributeKeyPricePrecision   = "price_precision"
	AttributeKeyLastExecutePrice = "last_execute_price"
	AttributeKeySender           = "sender"

	AttributeKeySequence    = "sequence"
	AttributeKeyOrderType   = "order_type"
	AttributeKeyPrice       = "price"
	AttributeKeyQuantity    = "quantity"
	AttributeKeySide        = "side"
	AttributeKeyHeight      = "height"
	AttributeKeyTimeInForce = "time_in_force"

	AttributeKeyDelOrderReason = "del_order_reason"
	AttributeKeyDelOrderHeight = "del_order_height"
	AttributeKeyUsedCommission = "used_commission"
	AttributeKeyLeftStock      = "left_stock"
	AttributeKeyDealStock      = "deal_stock"
	AttributeKeyRemainAmount   = "remain_amount"
	AttributeKeyDealMoney      = "deal_money"

	AttributeKeyEffectiveTime = "effective_time"

	AttributeKeyOldPricePrecision = "old_price_precision"
	AttributeKeyNewPricePrecision = "new_price_precision"
)
