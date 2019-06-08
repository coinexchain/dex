package market

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/types"
)

const (
	DefaultCreateMarketFee             = 1E12 // 10000 * 10 ^8
	DefaultFixedTradeFee               = 0
	DefaultGTEOrderLifetime            = 100
	DefaultMaxExecutedPriceChangeRatio = 25
	MarketFeeRatePrecision             = 4
	DefaultMarketFeeRate               = 0
)

var (
	KeyCreateMarketFee             = []byte("CreateMarketFee")
	KeyFixedTradeFee               = []byte("FixedTradeFee")
	KeyGTEOrderLifetime            = []byte("GTEOrderLifetime")
	KeyMaxExecutedPriceChangeRatio = []byte("MaxExecutedPriceChangeRatio")
	KeyMarketFeeRate               = []byte("MarketFeeRate")
)

type Params struct {
	CreateMarketFee             sdk.Coins `json:"create_market_fee"`
	FixedTradeFee               int64 `json:"fixed_trade_fee"`
	GTEOrderLifetime            int       `json:"gte_order_lifetime"`
	MaxExecutedPriceChangeRatio int       `json:"max_executed_price_change_ratio"`
	MarketFeeRate               int64     `json:"market_fee_rate"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyCreateMarketFee, Value: &p.CreateMarketFee},
		{Key: KeyFixedTradeFee, Value: &p.FixedTradeFee},
		{Key: KeyGTEOrderLifetime, Value: &p.GTEOrderLifetime},
		{Key: KeyMaxExecutedPriceChangeRatio, Value: &p.MaxExecutedPriceChangeRatio},
		{Key: KeyMarketFeeRate, Value: &p.MarketFeeRate},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		types.NewCetCoins(DefaultCreateMarketFee),
		DefaultFixedTradeFee,
		DefaultGTEOrderLifetime,
		DefaultMaxExecutedPriceChangeRatio,
		DefaultMarketFeeRate,
	}
}

func (p *Params) ValidateGenesis() error {
	if p.CreateMarketFee.Empty() || p.CreateMarketFee.IsAnyNegative() {
		return fmt.Errorf("%s must be a valid sdk.Coins, is %s", KeyCreateMarketFee, p.CreateMarketFee.String())
	}
	if p.FixedTradeFee<=0 {
		return fmt.Errorf("%s must be a valid sdk.Coins, is %d", KeyFixedTradeFee, p.FixedTradeFee)
	}

	if p.MaxExecutedPriceChangeRatio < 0 || p.MarketFeeRate < 0 || p.GTEOrderLifetime < 0 {
		return fmt.Errorf("params must be positive, MaxExecutedPriceChangeRatio "+
			": %d, MarketFeeRate: %d, GTEOrderLifetime : %d",
			p.MaxExecutedPriceChangeRatio, p.MarketFeeRate, p.GTEOrderLifetime)
	}
	return nil
}

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
