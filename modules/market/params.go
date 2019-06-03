package market

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/types"
)

const (
	CreateMarketFee             = 1E11 // 1000 * 10 ^8
	FilterStaleOrderInterval    = 10
	GTEOrderLifetime            = 100
	MaxExecutedPriceChangeRatio = 25
)

var (
	KeyCreateMarketFee             = []byte("CreateMarketFee")
	KeyFilterStaleOrderInterval    = []byte("FilterStaleOrderInterval")
	KeyGTEOrderLifetime            = []byte("GTEOrderLifetime")
	KeyMaxExecutedPriceChangeRatio = []byte("MaxExecutedPriceChangeRatio")
)

type Params struct {
	CreateMarketFee             sdk.Coins `json:"create_market_fee"`
	FilterStaleOrderInterval    int       `json:"filter_stale_order_interval"`
	GTEOrderLifetime            int       `json:"gte_order_lifetime"`
	MaxExecutedPriceChangeRatio int       `json:"max_executed_price_change_ratio"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyCreateMarketFee, Value: &p.CreateMarketFee},
		{Key: KeyFilterStaleOrderInterval, Value: &p.FilterStaleOrderInterval},
		{Key: KeyGTEOrderLifetime, Value: &p.GTEOrderLifetime},
		{Key: KeyMaxExecutedPriceChangeRatio, Value: &p.MaxExecutedPriceChangeRatio},
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
		types.NewCetCoins(CreateMarketFee),
		FilterStaleOrderInterval,
		GTEOrderLifetime,
		MaxExecutedPriceChangeRatio,
	}
}

func (p *Params) ValidateGenesis() error {
	if p.CreateMarketFee.Empty() || p.CreateMarketFee.IsAnyNegative() {
		return fmt.Errorf("%s must be a valid sdk.Coins, is %s", KeyCreateMarketFee, p.CreateMarketFee.String())
	}

	if p.MaxExecutedPriceChangeRatio < 0 || p.GTEOrderLifetime < 0 || p.FilterStaleOrderInterval < 0 {
		return fmt.Errorf("params must be positive, MaxExecutedPriceChangeRatio "+
			": %d, GTEOrderLifetime : %d, FilterStaleOrderInterval : %d ",
			p.MaxExecutedPriceChangeRatio, p.GTEOrderLifetime, p.MaxExecutedPriceChangeRatio)
	}
	return nil
}

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
