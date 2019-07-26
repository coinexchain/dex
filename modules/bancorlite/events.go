package bancorlite

// Market module event types
var (
	EventTypeBancorlite = "bancorlite"

	AttributeKeyCreateFor = "create-for"
	AttributeKeyTradeFor  = "trade-for"

	AttributeOwner          = "bancor-owner"
	AttributeMaxSupply      = "bancor-max-supply"
	AttributeNewStockInPool = "bancor-new-stock-in-pool"
	AttributeNewMoneyInPool = "bancor-new-money-in-pool"
	AttributeNewPrice       = "bancor-new-price"
	AttributeCoinsFromPool  = "bancor-coins-from-pool"
	AttributeCoinsToPool    = "bancor-coins-to-pool"
	AttributeTradeSide      = "bancor-trade-side"

	KafkaBancorTrade  = "bancor_trade"
	KafkaBancorCreate = "bancor_create"
	KafkaBancorCancel = "bancor_cancel"
	KafkaBancorInfo   = "bancor_info"
)
