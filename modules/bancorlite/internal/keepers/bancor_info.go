package keepers

import (
	"fmt"
	"time"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dex "github.com/coinexchain/dex/types"
)

type BancorInfo struct {
	Owner              sdk.AccAddress `json:"sender"`
	Stock              string         `json:"stock"`
	Money              string         `json:"money"`
	InitPrice          sdk.Dec        `json:"init_price"`
	MaxSupply          sdk.Int        `json:"max_supply"`
	StockPrecision     byte           `json:"stock_precision"`
	MaxPrice           sdk.Dec        `json:"max_price"`
	MaxMoney           sdk.Int        `json:"max_money"`
	AR                 int64          `json:"ar"`
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
	if bi.MaxMoney.IsZero() {
		bi.Price = bi.MaxPrice.Sub(bi.InitPrice).MulInt(suppliedStock).QuoInt(bi.MaxSupply).Add(bi.InitPrice)
		bi.MoneyInPool = bi.Price.Add(bi.InitPrice).MulInt(suppliedStock).QuoInt64(2).RoundInt()
	} else {
		// s = s/s_max * 1000, as of precision is 0.001
		factoredStock := suppliedStock.MulRaw(types.SupplyRatioSamples)
		s := factoredStock.Quo(bi.MaxSupply).Int64()
		if s > types.SupplyRatioSamples {
			return false
		}
		contrast := sdk.NewInt(s).Mul(bi.MaxSupply)
		// ratio = (s/s_max)^ar, ar = (p_max * s_max - m_max) / (m_max - p_init * s_max)
		ratio := types.TableLookup(bi.AR+10, s)
		if contrast.GT(factoredStock) {
			if s == types.SupplyRatioSamples {
				return false
			}
			ratioNear := types.TableLookup(bi.AR+10, s-1)
			// ratio = (ratio - ratioNear) * (stock_now / s_max * 1000 - (s-1)) + ratioNear
			ratio = ratio.Sub(ratioNear).MulInt(factoredStock.Sub(sdk.NewInt(s - 1).Mul(bi.MaxSupply))).
				Quo(sdk.NewDecFromInt(bi.MaxSupply)).Add(ratioNear)
		} else if factoredStock.GT(contrast) {
			if s > types.SupplyRatioSamples {
				return false
			}
			// ratio = (ratioNear - ratio) * (stock_now / s_max * 1000 - (s)) + ratio
			ratioNear := types.TableLookup(bi.AR+10, s+1)
			ratio = ratioNear.Sub(ratio).MulInt(factoredStock.Sub(sdk.NewInt(s).Mul(bi.MaxSupply))).
				Quo(sdk.NewDecFromInt(bi.MaxSupply)).Add(ratio)
		}

		// m_now = m_max * ratio
		bi.MoneyInPool = ratio.MulInt(bi.MaxMoney).TruncateInt()
		// price_ratio = (s/s_max)^(ar)
		priceRatio := types.TableLookup(bi.AR, s)
		// price = priceRatio * (maxPrice - initPrice) + initPrice
		bi.Price = priceRatio.MulTruncate(bi.MaxPrice.Sub(bi.InitPrice)).Add(bi.InitPrice)

	}
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
	StockPrecision     string `json:"stock_precision"`
	MaxPrice           string `json:"max_price"`
	MaxMoney           string `json:"max_money"`
	AR                 string `json:"ar"`
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
		StockPrecision:     fmt.Sprintf("%d", bi.StockPrecision),
		MaxPrice:           bi.MaxPrice.String(),
		MaxMoney:           bi.MaxMoney.String(),
		AR:                 fmt.Sprintf("%d", bi.AR),
		CurrentPrice:       bi.Price.String(),
		StockInPool:        bi.StockInPool.String(),
		MoneyInPool:        bi.MoneyInPool.String(),
		EarliestCancelTime: time.Unix(bi.EarliestCancelTime, 0).Format(time.RFC3339),
	}
}
