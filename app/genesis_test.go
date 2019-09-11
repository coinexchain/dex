package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/market"
)

func TestFromToMap(t *testing.T) {
	gsMap := ModuleBasics.DefaultGenesis()
	cdc := MakeCodec()
	m := FromMap(cdc, gsMap)
	m.toMap(cdc)
}

func TestDefaultGenesisState(t *testing.T) {
	state := ModuleBasics.DefaultGenesis()

	//auth
	var authData auth.GenesisState
	auth.ModuleCdc.MustUnmarshalJSON(state[auth.ModuleName], &authData)
	require.Equal(t, 512, int(authData.Params.MaxMemoCharacters))

	// staking
	var stakingData staking.GenesisState
	staking.ModuleCdc.MustUnmarshalJSON(state[staking.ModuleName], &stakingData)
	require.Equal(t, "cet", stakingData.Params.BondDenom)
	require.Equal(t, "504h0m0s", stakingData.Params.UnbondingTime.String()) // 21 days
	require.Equal(t, 42, int(stakingData.Params.MaxValidators))
	require.Equal(t, 7, int(stakingData.Params.MaxEntries))

	// slashing
	var slashingData slashing.GenesisState
	slashing.ModuleCdc.MustUnmarshalJSON(state[slashing.ModuleName], &slashingData)
	require.Equal(t, "504h0m0s", slashingData.Params.MaxEvidenceAge.String())
	require.Equal(t, "10m0s", slashingData.Params.DowntimeJailDuration.String())
	require.Equal(t, 10000, int(slashingData.Params.SignedBlocksWindow))
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), slashingData.Params.MinSignedPerWindow)
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), slashingData.Params.SlashFractionDoubleSign)
	require.Equal(t, sdk.MustNewDecFromStr("0.0001"), slashingData.Params.SlashFractionDowntime)

	// distr
	var distrData distr.GenesisState
	distr.ModuleCdc.MustUnmarshalJSON(state[distr.ModuleName], &distrData)
	require.True(t, distrData.WithdrawAddrEnabled)
	require.Equal(t, sdk.MustNewDecFromStr("0.02"), distrData.CommunityTax)
	require.Equal(t, sdk.MustNewDecFromStr("0.01"), distrData.BaseProposerReward)
	require.Equal(t, sdk.MustNewDecFromStr("0.04"), distrData.BonusProposerReward)

	// gov
	var govData gov.GenesisState
	gov.ModuleCdc.MustUnmarshalJSON(state[gov.ModuleName], &govData)
	require.Equal(t, "1000000000000cet", govData.DepositParams.MinDeposit.String())
	require.Equal(t, "336h0m0s", govData.DepositParams.MaxDepositPeriod.String())
	require.Equal(t, "336h0m0s", govData.VotingParams.VotingPeriod.String())
	require.Equal(t, sdk.MustNewDecFromStr("0.4"), govData.TallyParams.Quorum)
	require.Equal(t, sdk.MustNewDecFromStr("0.5"), govData.TallyParams.Threshold)
	require.Equal(t, sdk.MustNewDecFromStr("0.334"), govData.TallyParams.Veto)

	// crisis
	var crisisData crisis.GenesisState
	crisis.ModuleCdc.MustUnmarshalJSON(state[crisis.ModuleName], &crisisData)
	require.Equal(t, "35000000000000cet", crisisData.ConstantFee.String())

	// others
	var authxData authx.GenesisState
	authx.ModuleCdc.MustUnmarshalJSON(state[authx.ModuleName], &authxData)
	require.Equal(t, sdk.NewDec(20), authxData.Params.MinGasPriceLimit)

	// market
	var marketData market.GenesisState
	market.ModuleCdc.MustUnmarshalJSON(state[market.ModuleName], &marketData)
	require.Equal(t, 0, len(marketData.MarketInfos))
	require.Equal(t, 0, len(marketData.Orders))
	require.Equal(t, int64(0), marketData.OrderCleanTime)
	require.Equal(t, 10000, marketData.Params.GTEOrderLifetime)
	require.Equal(t, int64(6000000), marketData.Params.GTEOrderFeatureFeeByBlocks)
	require.Equal(t, int64(1e12), marketData.Params.CreateMarketFee)
}
