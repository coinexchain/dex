package market

type GenesisStateOfMarket struct {
	Params      ParamsOfMarket `json:"params"`
	Orders      []Order        `json:"orders"`
	MarketInfos []MarketInfo   `json:"market_infos"`
}

func NewGenesisStateOfMarket(params ParamsOfMarket, orders []Order, infos []MarketInfo) GenesisStateOfMarket {
	return GenesisStateOfMarket{
		Params:      params,
		Orders:      orders,
		MarketInfos: infos,
	}
}
