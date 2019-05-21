package market

type MarketInfo struct {
	Stock             string
	Money             string
	Create            string
	PricePrecision    byte
	LastExecutedPrice int
}

func (minfo *MarketInfo) CheckMarketInfoValid() bool {
	return true
}
