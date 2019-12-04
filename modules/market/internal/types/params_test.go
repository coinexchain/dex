package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	params := Params{
		CreateMarketFee:             100,
		MarketMinExpiredTime:        100,
		GTEOrderLifetime:            100,
		GTEOrderFeatureFeeByBlocks:  100,
		MaxExecutedPriceChangeRatio: 100,
		MarketFeeRate:               100,
		MarketFeeMin:                100,
		FeeForZeroDeal:              100,
	}
	require.Equal(t, nil, params.ValidateGenesis())
	params1 := params
	params1.CreateMarketFee = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.GTEOrderLifetime = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.GTEOrderFeatureFeeByBlocks = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.MaxExecutedPriceChangeRatio = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.MarketFeeRate = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.MarketFeeMin = -1
	require.NotNil(t, params1.ValidateGenesis())
	params1 = params
	params1.FeeForZeroDeal = -1
	require.NotNil(t, params1.ValidateGenesis())
}
