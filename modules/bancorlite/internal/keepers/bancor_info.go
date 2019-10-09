package keepers

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dex "github.com/coinexchain/dex/types"
)

type BancorInfo struct {
	Owner              sdk.AccAddress `json:"sender"`
	Stock              string         `json:"stock"`
	Money              string         `json:"money"`
	InitPrice          sdk.Dec        `json:"init_price"`
	MaxSupply          sdk.Int        `json:"max_supply"`
	MaxPrice           sdk.Dec        `json:"max_price"`
	Price              sdk.Dec        `json:"price"`
	StockInPool        sdk.Int        `json:"stock_in_pool"`
	MoneyInPool        sdk.Int        `json:"money_in_pool"`
	EarliestCancelTime int64          `json:"earliest_cancel_time"`
}

func (bi *BancorInfo) GetSymbol() string {
	return dex.GetSymbol(bi.Stock, bi.Money)
}

func (bi *BancorInfo) UpdateStockInPool(stockInPool sdk.Int) bool {
	if stockInPool.IsNegative() || stockInPool.GT(bi.MaxSupply) {
		return false
	}

	bi.StockInPool = stockInPool
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	bi.Price = bi.MaxPrice.Sub(bi.InitPrice).MulInt(suppliedStock).QuoInt(bi.MaxSupply).Add(bi.InitPrice)
	bi.MoneyInPool = bi.Price.Add(bi.InitPrice).MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return true
}

func (bi *BancorInfo) IsConsistent() bool {
	if bi.StockInPool.IsNegative() || bi.StockInPool.GT(bi.MaxSupply) {
		return false
	}
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	price := bi.MaxPrice.Sub(bi.InitPrice).MulInt(suppliedStock).QuoInt(bi.MaxSupply).Add(bi.InitPrice)
	moneyInPool := price.Add(bi.InitPrice).MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return price.Equal(bi.Price) && moneyInPool.Equal(bi.MoneyInPool)
}

type BancorInfoDisplay struct {
	Stock              string `json:"stock"`
	Money              string `json:"money"`
	InitPrice          string `json:"init_price"`
	MaxSupply          string `json:"max_supply"`
	MaxPrice           string `json:"max_price"`
	CurrentPrice       string `json:"current_price"`
	StockInPool        string `json:"stock_in_pool"`
	MoneyInPool        string `json:"money_in_pool"`
	EarliestCancelTime string `json:"earliest_cancel_time"`
}

func NewBancorInfoDisplay(bi *BancorInfo) BancorInfoDisplay {
	return BancorInfoDisplay{
		Stock:              bi.Stock,
		Money:              bi.Money,
		InitPrice:          bi.InitPrice.String(),
		MaxSupply:          bi.MaxSupply.String(),
		MaxPrice:           bi.MaxPrice.String(),
		CurrentPrice:       bi.Price.String(),
		StockInPool:        bi.StockInPool.String(),
		MoneyInPool:        bi.MoneyInPool.String(),
		EarliestCancelTime: time.Unix(bi.EarliestCancelTime, 0).Format(time.RFC3339),
	}
}
