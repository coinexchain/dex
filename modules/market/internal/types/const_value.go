package types

type OrderType = byte

const (
	MinTokenPricePrecision           = 0
	MaxTokenPricePrecision           = 18
	LimitOrder             OrderType = 2
	SymbolSeparator                  = "/"
	OrderIDSeparator                 = "-"
	ExtraFrozenMoney                 = 0 // 100
	OrderIDPartsNum                  = 2
)

const (
	BID = 1
	BUY = 1

	ASK  = 2
	SELL = 2
)

const (
	DecByteCount = 40 // Dec's BitLen would not be larger than 255+60, so 40 bytes are enough
	GTE          = 3
	IOC          = 4
	LIMIT        = 2
)

const (
	IntegrationNetSubString       = "coinex-integrationtest"
	MaxOrderAmount          int64 = 1e18
)
