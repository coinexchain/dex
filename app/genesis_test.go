package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.NoError(t, state.Validate())

	// staking
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "504h0m0s", state.StakingData.Params.UnbondingTime.String()) // 21 days
	require.Equal(t, 42, int(state.StakingData.Params.MaxValidators))
	require.Equal(t, 7, int(state.StakingData.Params.MaxEntries))

	// slashing
	require.Equal(t, "504h0m0s", state.SlashingData.Params.MaxEvidenceAge.String())
	require.Equal(t, "10m0s", state.SlashingData.Params.DowntimeJailDuration.String())
	require.Equal(t, 1000, int(state.SlashingData.Params.SignedBlocksWindow))
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), state.SlashingData.Params.MinSignedPerWindow)
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), state.SlashingData.Params.SlashFractionDoubleSign)
	require.Equal(t, sdk.MustNewDecFromStr("0.0001"), state.SlashingData.Params.SlashFractionDowntime)

	// distr
	require.True(t, state.DistrData.WithdrawAddrEnabled)
	require.Equal(t, sdk.MustNewDecFromStr("0.02"), state.DistrData.CommunityTax)
	require.Equal(t, sdk.MustNewDecFromStr("0.01"), state.DistrData.BaseProposerReward)
	require.Equal(t, sdk.MustNewDecFromStr("0.04"), state.DistrData.BonusProposerReward)

	// gov
	require.Equal(t, "10000000cet", state.GovData.DepositParams.MinDeposit.String())
	require.Equal(t, "336h0m0s", state.GovData.DepositParams.MaxDepositPeriod.String())
	require.Equal(t, "336h0m0s", state.GovData.VotingParams.VotingPeriod.String())
	require.Equal(t, sdk.MustNewDecFromStr("0.4"), state.GovData.TallyParams.Quorum)
	require.Equal(t, sdk.MustNewDecFromStr("0.5"), state.GovData.TallyParams.Threshold)
	require.Equal(t, sdk.MustNewDecFromStr("0.334"), state.GovData.TallyParams.Veto)

	// crisis
	require.Equal(t, "1000cet", state.CrisisData.ConstantFee.String())

	// others
	require.Equal(t, sdk.NewDec(20), state.AuthXData.Params.MinGasPriceLimit)
}
