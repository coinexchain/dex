package market

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultCreateMarketFee             = 1E12 // 10000 * 10 ^8
	DefaultFixedTradeFee               = 1000000
	DefaultGTEOrderLifetime            = 10000
	DefaultGTEOrderFeatureFeeByBlocks  = 6000000
	DefaultMaxExecutedPriceChangeRatio = 25
	MarketFeeRatePrecision             = 4
	DefaultMarketFeeRate               = 10
	DefaultMarketFeeMin                = 1000000
	DefaultFeeForZeroDeal              = 1000000
	DefaultMarketMinExpiredTime        = 60 * 60 * 24 * 7
	DefaultChainIDVersion              = 0
)

var (
	KeyChainIDVersion              = []byte("ChainIDVersion")
	KeyCreateMarketFee             = []byte("CreateMarketFee")
	KeyFixedTradeFee               = []byte("FixedTradeFee")
	keyMarketMinExpiredTime        = []byte("MarketMinExpiredTime")
	KeyGTEOrderLifetime            = []byte("GTEOrderLifetime")
	KeyGTEOrderFeatureFeeByBlocks  = []byte("GTEOrderFeatureFeeByBlocks")
	KeyMaxExecutedPriceChangeRatio = []byte("MaxExecutedPriceChangeRatio")
	KeyMarketFeeRate               = []byte("MarketFeeRate")
	KeyMarketFeeMin                = []byte("MarketFeeMin")
	KeyFeeForZeroDeal              = []byte("FeeForZeroDeal")
)

type Params struct {
	CreateMarketFee             int64 `json:"create_market_fee"`
	FixedTradeFee               int64 `json:"fixed_trade_fee"`
	MarketMinExpiredTime        int64 `json:"market_min_expired_time"`
	GTEOrderLifetime            int   `json:"gte_order_lifetime"`
	GTEOrderFeatureFeeByBlocks  int64 `json:"gte_order_feature_fee_by_blocks"`
	MaxExecutedPriceChangeRatio int   `json:"max_executed_price_change_ratio"`
	MarketFeeRate               int64 `json:"market_fee_rate"`
	MarketFeeMin                int64 `json:"market_fee_min"`
	FeeForZeroDeal              int64 `json:"fee_for_zero_deal"`
	ChainIDVersion              int64 `json:"chain_id_version"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyCreateMarketFee, Value: &p.CreateMarketFee},
		{Key: KeyFixedTradeFee, Value: &p.FixedTradeFee},
		{Key: keyMarketMinExpiredTime, Value: &p.MarketMinExpiredTime},
		{Key: KeyGTEOrderLifetime, Value: &p.GTEOrderLifetime},
		{Key: KeyGTEOrderFeatureFeeByBlocks, Value: &p.GTEOrderFeatureFeeByBlocks},
		{Key: KeyMaxExecutedPriceChangeRatio, Value: &p.MaxExecutedPriceChangeRatio},
		{Key: KeyMarketFeeRate, Value: &p.MarketFeeRate},
		{Key: KeyMarketFeeMin, Value: &p.MarketFeeMin},
		{Key: KeyFeeForZeroDeal, Value: &p.FeeForZeroDeal},
		{Key: KeyChainIDVersion, Value: &p.ChainIDVersion},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		DefaultCreateMarketFee,
		DefaultFixedTradeFee,
		DefaultMarketMinExpiredTime,
		DefaultGTEOrderLifetime,
		DefaultGTEOrderFeatureFeeByBlocks,
		DefaultMaxExecutedPriceChangeRatio,
		DefaultMarketFeeRate,
		DefaultMarketFeeMin,
		DefaultFeeForZeroDeal,
		DefaultChainIDVersion,
	}
}

func (p *Params) ValidateGenesis() error {
	if p.CreateMarketFee <= 0 {
		return fmt.Errorf("%s must be a positive number, is %d", KeyCreateMarketFee, p.CreateMarketFee)
	}
	if p.FixedTradeFee < 0 {
		return fmt.Errorf("%s must be a valid sdk.Coins, is %d", KeyFixedTradeFee, p.FixedTradeFee)
	}

	if p.MaxExecutedPriceChangeRatio < 0 || p.MarketFeeRate < 0 || p.MarketFeeMin < 0 || p.FeeForZeroDeal < 0 || p.GTEOrderLifetime < 0 {
		return fmt.Errorf("params must be positive, MaxExecutedPriceChangeRatio "+
			": %d, MarketFeeRate: %d, MarketFeeMin: %d, FeeForZeroDeal: %d, GTEOrderLifetime : %d",
			p.MaxExecutedPriceChangeRatio, p.MarketFeeRate, p.MarketFeeMin, p.FeeForZeroDeal, p.GTEOrderLifetime)
	}
	return nil
}

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
