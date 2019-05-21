package market

type Order struct {
	Sender         string
	Sequence       uint64
	Symbol         string
	OrderType      byte
	PricePrecision byte
	Price          uint64
	Quantity       uint64
	Side           byte
	TimeInForce    int
	Height         uint64

	// These field will change when order filled/cancel.
	LeftStock uint64
	Freeze    uint64
	DealStock uint64
	DealMoney uint64
}

func (or *Order) CheckOrderValidToMempool() bool {

	return true
}

func (or *Order) CheckOrderValidToExecute() bool {

	return true
}
