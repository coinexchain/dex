package bancorlite

// Market module event types
var (
	AttributeValueCategory = ModuleName

	EventTypeKeyBancorInit   = "bancor_init"
	EventTypeKeyBancorTrade  = "bancor_trade"
	EventTypeKeyBancorCancel = "bancor_cancel"

	AttributeSymbol         = "symbol"
	AttributeOwner          = "bancor_owner"
	AttributeMaxSupply      = "bancor_max_supply"
	AttributeNewStockInPool = "bancor_new_stock_in_pool"
	AttributeNewMoneyInPool = "bancor_new_money_in_pool"
	AttributeNewPrice       = "bancor_new_price"
	AttributeCoinsFromPool  = "bancor_coins_from_pool"
	AttributeCoinsToPool    = "bancor_coins_to_pool"
	AttributeTradeSide      = "bancor_trade_side"

	KafkaBancorTrade  = "bancor_trade"
	KafkaBancorCreate = "bancor_create"
	KafkaBancorCancel = "bancor_cancel"
	KafkaBancorInfo   = "bancor_info"
)
