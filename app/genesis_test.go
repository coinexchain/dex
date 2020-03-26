package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/cet-sdk/modules/alias"
	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/bancorlite"
	"github.com/coinexchain/cet-sdk/modules/bankx"
	"github.com/coinexchain/cet-sdk/modules/incentive"
	"github.com/coinexchain/cet-sdk/modules/market"
	"github.com/coinexchain/cet-sdk/modules/stakingx"
)

func TestFromToMap(t *testing.T) {
	gsMap := ModuleBasics.DefaultGenesis()
	cdc := MakeCodec()
	m := FromMap(cdc, gsMap)
	m.toMap(cdc)
}

func TestDefaultGenesisState(t *testing.T) {
	state := ModuleBasics.DefaultGenesis()

	// auth
	var authData auth.GenesisState
	auth.ModuleCdc.MustUnmarshalJSON(state[auth.ModuleName], &authData)
	require.Equal(t, uint64(512), authData.Params.MaxMemoCharacters)
	require.Equal(t, uint64(7), authData.Params.TxSigLimit)
	require.Equal(t, DefaultTxSizeCostPerByte, authData.Params.TxSizeCostPerByte)
	require.Equal(t, DefaultSigVerifyCostED25519, authData.Params.SigVerifyCostED25519)
	require.Equal(t, DefaultSigVerifyCostSecp256k1, authData.Params.SigVerifyCostSecp256k1)

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
	checkCET(t, 10000, govData.DepositParams.MinDeposit)
	require.Equal(t, "336h0m0s", govData.DepositParams.MaxDepositPeriod.String())
	require.Equal(t, "168h0m0s", govData.VotingParams.VotingPeriod.String())
	require.Equal(t, sdk.MustNewDecFromStr("0.4"), govData.TallyParams.Quorum)
	require.Equal(t, sdk.MustNewDecFromStr("0.5"), govData.TallyParams.Threshold)
	require.Equal(t, sdk.MustNewDecFromStr("0.334"), govData.TallyParams.Veto)

	// crisis
	var crisisData crisis.GenesisState
	crisis.ModuleCdc.MustUnmarshalJSON(state[crisis.ModuleName], &crisisData)
	checkCET(t, 100000, sdk.Coins{crisisData.ConstantFee})

	// others
	var authxData authx.GenesisState
	authx.ModuleCdc.MustUnmarshalJSON(state[authx.ModuleName], &authxData)
	require.Equal(t, sdk.NewDec(20), authxData.Params.MinGasPriceLimit)
	var bankxData bankx.GenesisState
	bankx.ModuleCdc.MustUnmarshalJSON(state[bankx.ModuleName], &bankxData)
	require.Equal(t, int64(1e8), bankxData.Params.ActivationFee)
	require.Equal(t, int64(604800e9), bankxData.Params.LockCoinsFreeTime)
	require.Equal(t, int64(1000000), bankxData.Params.LockCoinsFeePerDay)
	var stakingxData stakingx.GenesisState
	bankx.ModuleCdc.MustUnmarshalJSON(state[stakingx.ModuleName], &stakingxData) // TODO
	require.Equal(t, int64(5000000e8), stakingxData.Params.MinSelfDelegation)
	require.Equal(t, sdk.MustNewDecFromStr("0.1"), stakingxData.Params.MinMandatoryCommissionRate)

	// alias
	var aliasData alias.GenesisState
	alias.ModuleCdc.MustUnmarshalJSON(state[alias.ModuleName], &aliasData)
	require.Equal(t, 5, aliasData.Params.MaxAliasCount)
	require.Equal(t, int64(10000e8), aliasData.Params.FeeForAliasLength2)
	require.Equal(t, int64(5000e8), aliasData.Params.FeeForAliasLength3)
	require.Equal(t, int64(2000e8), aliasData.Params.FeeForAliasLength4)
	require.Equal(t, int64(1000e8), aliasData.Params.FeeForAliasLength5)
	require.Equal(t, int64(100e8), aliasData.Params.FeeForAliasLength6)
	require.Equal(t, int64(10e8), aliasData.Params.FeeForAliasLength7OrHigher)

	// asset
	var assetData asset.GenesisState
	asset.ModuleCdc.MustUnmarshalJSON(state[asset.ModuleName], &assetData)
	require.Equal(t, int64(10000e8), assetData.Params.IssueRareTokenFee)
	require.Equal(t, int64(1000e8), assetData.Params.Issue3CharTokenFee)
	require.Equal(t, int64(500e8), assetData.Params.Issue4CharTokenFee)
	require.Equal(t, int64(200e8), assetData.Params.Issue5CharTokenFee)
	require.Equal(t, int64(100e8), assetData.Params.Issue6CharTokenFee)
	require.Equal(t, int64(50e8), assetData.Params.IssueTokenFee)

	// bancor
	var bancorData bancorlite.GenesisState
	bancorlite.ModuleCdc.MustUnmarshalJSON(state[bancorlite.ModuleName], &bancorData)
	require.Equal(t, int64(100e8), bancorData.Params.CreateBancorFee)
	require.Equal(t, int64(100e8), bancorData.Params.CancelBancorFee)
	require.Equal(t, int64(10), bancorData.Params.TradeFeeRate)

	// incentive
	var incentiveData incentive.GenesisState
	incentive.ModuleCdc.MustUnmarshalJSON(state[incentive.ModuleName], &incentiveData)
	require.Equal(t, int64(2e8), incentiveData.Params.DefaultRewardPerBlock)

	// market
	var marketData market.GenesisState
	market.ModuleCdc.MustUnmarshalJSON(state[market.ModuleName], &marketData)
	require.Equal(t, int64(100e8), marketData.Params.CreateMarketFee)
	require.Equal(t, int64(604800e9), marketData.Params.MarketMinExpiredTime)
	require.EqualValues(t, 200000, marketData.Params.GTEOrderLifetime)
	require.Equal(t, int64(10), marketData.Params.GTEOrderFeatureFeeByBlocks)
	require.EqualValues(t, 25, marketData.Params.MaxExecutedPriceChangeRatio)
	require.Equal(t, int64(10), marketData.Params.MarketFeeRate)
	require.Equal(t, int64(1000000), marketData.Params.MarketFeeMin)
	require.Equal(t, int64(1000000), marketData.Params.FeeForZeroDeal)
}

func checkCET(t *testing.T, amt int64, coins sdk.Coins) {
	require.Equal(t, fmt.Sprintf("%dcet", amt*1e8), coins.String())
}
