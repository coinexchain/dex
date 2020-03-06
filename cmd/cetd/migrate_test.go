package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/cet-sdk/modules/bancorlite"
	"github.com/coinexchain/cet-sdk/modules/market"
	"github.com/coinexchain/dex/app"
)

func TestMigrate(t *testing.T) {
	state := app.NewDefaultGenesisState()

	// simulate DEX1
	state.AssetData.Params.Issue4CharTokenFee = 0
	state.AuthXData.Params.RebateRatio = 0
	state.MarketData.Orders = append(state.MarketData.Orders,
		&market.Order{FrozenFee: 100})
	state.BancorData.BancorInfoMap["x"] = bancorlite.BancorInfo{}

	// upgrade to DEX2
	upgradeGenesisState(&state)

	// check state
	require.Equal(t, app.VotingPeriod, state.GovData.VotingParams.VotingPeriod)
	require.Equal(t, int64(1000), state.AuthXData.Params.RebateRatio)
	require.Equal(t, int64(5e11), state.AssetData.Params.Issue4CharTokenFee)
	require.Equal(t, int64(100), state.MarketData.Orders[0].FrozenCommission)
	require.Equal(t, sdk.ZeroInt(), state.BancorData.BancorInfoMap["x"].MaxMoney)
}
