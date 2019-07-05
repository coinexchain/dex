package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

func TestDefaultGenesisState(t *testing.T) {
	state := NewDefaultGenesisState()
	require.NoError(t, state.Validate())

	//auth
	require.Equal(t, 512, int(state.AuthData.Params.MaxMemoCharacters))

	// staking
	require.Equal(t, "cet", state.StakingData.Params.BondDenom)
	require.Equal(t, "504h0m0s", state.StakingData.Params.UnbondingTime.String()) // 21 days
	require.Equal(t, 42, int(state.StakingData.Params.MaxValidators))
	require.Equal(t, 7, int(state.StakingData.Params.MaxEntries))

	// slashing
	require.Equal(t, "504h0m0s", state.SlashingData.Params.MaxEvidenceAge.String())
	require.Equal(t, "10m0s", state.SlashingData.Params.DowntimeJailDuration.String())
	require.Equal(t, 10000, int(state.SlashingData.Params.SignedBlocksWindow))
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), state.SlashingData.Params.MinSignedPerWindow)
	require.Equal(t, sdk.MustNewDecFromStr("0.05"), state.SlashingData.Params.SlashFractionDoubleSign)
	require.Equal(t, sdk.MustNewDecFromStr("0.0001"), state.SlashingData.Params.SlashFractionDowntime)

	// distr
	require.True(t, state.DistrData.WithdrawAddrEnabled)
	require.Equal(t, sdk.MustNewDecFromStr("0.02"), state.DistrData.CommunityTax)
	require.Equal(t, sdk.MustNewDecFromStr("0.01"), state.DistrData.BaseProposerReward)
	require.Equal(t, sdk.MustNewDecFromStr("0.04"), state.DistrData.BonusProposerReward)

	// gov
	require.Equal(t, "1000000000000cet", state.GovData.DepositParams.MinDeposit.String())
	require.Equal(t, "336h0m0s", state.GovData.DepositParams.MaxDepositPeriod.String())
	require.Equal(t, "336h0m0s", state.GovData.VotingParams.VotingPeriod.String())
	require.Equal(t, sdk.MustNewDecFromStr("0.4"), state.GovData.TallyParams.Quorum)
	require.Equal(t, sdk.MustNewDecFromStr("0.5"), state.GovData.TallyParams.Threshold)
	require.Equal(t, sdk.MustNewDecFromStr("0.334"), state.GovData.TallyParams.Veto)

	// crisis
	require.Equal(t, "35000000000000cet", state.CrisisData.ConstantFee.String())

	// others
	require.Equal(t, sdk.NewDec(20), state.AuthXData.Params.MinGasPriceLimit)
}

func TestGenesisAccountToAccount(t *testing.T) {

	_, _, addr := testutil.KeyPubAddr()
	bAcc := auth.NewBaseAccountWithAddress(addr)
	bAcc.SetCoins(types.NewCetCoins(1000))
	continuousAcc := auth.NewContinuousVestingAccount(&bAcc, 10000, 12345)
	delayedAcc := auth.NewDelayedVestingAccount(&bAcc, 12345)

	gcAcc := NewGenesisAccountI(continuousAcc)
	gdAcc := NewGenesisAccountI(delayedAcc)

	toAcc1 := gcAcc.ToAccount()
	toAcc2 := gdAcc.ToAccount()

	require.Equal(t, continuousAcc.String(), toAcc1.String())
	require.Equal(t, delayedAcc.String(), toAcc2.String())

}

func TestGenesisAccountValidate(t *testing.T) {

	duplicatedAccounts := []GenesisAccount{
		{Address: sdk.AccAddress("myaddr"), Coins: types.NewCetCoins(100), OriginalVesting: types.NewCetCoins(100)},
		{Address: sdk.AccAddress("myaddr"), Coins: types.NewCetCoins(100), OriginalVesting: types.NewCetCoins(200)},
	}
	invalidDelayedVestingAccounts := []GenesisAccount{
		{Address: sdk.AccAddress("myaddr1"), Coins: types.NewCetCoins(100), OriginalVesting: types.NewCetCoins(100)},
	}
	invalidContinuousAccounts := []GenesisAccount{
		{Address: sdk.AccAddress("myaddr2"), Coins: types.NewCetCoins(100), OriginalVesting: types.NewCetCoins(100), StartTime: 10000, EndTime: 900},
	}

	require.NotNil(t, validateGenesisStateAccounts(duplicatedAccounts))
	require.NotNil(t, validateGenesisStateAccounts(invalidContinuousAccounts))
	require.NotNil(t, validateGenesisStateAccounts(invalidDelayedVestingAccounts))

}
