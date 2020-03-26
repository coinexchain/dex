package main

import (
	"testing"
	"time"

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
	require.Equal(t, time.Hour*24*7, state.GovData.VotingParams.VotingPeriod)
	require.EqualValues(t, time.Hour*24*7, state.AuthXData.Params.RefereeChangeMinInterval)
	require.EqualValues(t, 2000, state.AuthXData.Params.RebateRatio)
	require.EqualValues(t, 1e14, state.StakingXData.Params.MinSelfDelegation)
	require.EqualValues(t, 1e12, state.AssetData.Params.IssueRareTokenFee)
	require.EqualValues(t, 1e11, state.AssetData.Params.Issue3CharTokenFee)
	require.EqualValues(t, 5e10, state.AssetData.Params.Issue4CharTokenFee)
	require.EqualValues(t, 2e10, state.AssetData.Params.Issue5CharTokenFee)
	require.EqualValues(t, 1e10, state.AssetData.Params.Issue6CharTokenFee)
	require.EqualValues(t, 5e9, state.AssetData.Params.IssueTokenFee)
	require.EqualValues(t, 1e10, state.MarketData.Params.CreateMarketFee)
	require.EqualValues(t, 200000, state.MarketData.Params.GTEOrderLifetime)
	require.EqualValues(t, 100, state.MarketData.Orders[0].FrozenCommission)
	require.Equal(t, sdk.ZeroInt(), state.BancorData.BancorInfoMap["x"].MaxMoney)
}
