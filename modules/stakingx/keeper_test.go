package stakingx

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func TestInitExportGenesis(t *testing.T) {
	sxk, ctx, _ := setUpInput()

	genesisState := GenesisState{
		Params: Params{
			MinSelfDelegation:          sdk.NewInt(DefaultMinSelfDelegation),
			MinMandatoryCommissionRate: DefaultMinMandatoryCommissionRate,
		},
	}

	InitGenesis(ctx, sxk, genesisState)
	exportGenesis := ExportGenesis(ctx, sxk)
	require.Equal(t, genesisState, exportGenesis)
}

func TestCalcBondPoolStatus(t *testing.T) {
	//initialize keeper & params &state
	sxk, ctx, _ := setUpInput()

	_, _, addr := testutil.KeyPubAddr()
	testParam := Params{
		MinSelfDelegation: sdk.ZeroInt(),
	}
	acc := auth.BaseAccount{
		Address: addr,
		Coins:   dex.NewCetCoins(1e8),
	}
	vacc := auth.NewDelayedVestingAccount(&acc, math.MaxInt64)
	sxk.ak.SetAccount(ctx, vacc)
	InitGenesis(ctx, sxk, GenesisState{Params: testParam})

	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(1000)),
	}
	sxk.dk.SetFeePool(ctx, feePool)

	bondedAcc := sxk.supplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName)
	bondedAcc.SetCoins(dex.NewCetCoins(1000))
	sxk.ak.SetAccount(ctx, bondedAcc)

	bondedAcc = sxk.supplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName)
	notBondedAcc := sxk.supplyKeeper.GetModuleAccount(ctx, staking.NotBondedPoolName)
	expectedNonBondableTokens := feePool.CommunityPool.AmountOf("cet").Add(acc.Coins.AmountOf("cet").ToDec())
	expectedTotalSupply := sxk.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("cet")
	expectedBondRatio := bondedAcc.GetCoins().AmountOf("cet").ToDec().QuoInt(expectedTotalSupply.Sub(expectedNonBondableTokens.RoundInt()))

	//test
	bondPool := sxk.CalcBondPoolStatus(ctx)

	require.Equal(t, expectedNonBondableTokens, bondPool.NonBondableTokens.ToDec())
	require.Equal(t, expectedTotalSupply, bondPool.TotalSupply)
	require.Equal(t, expectedBondRatio, bondPool.BondRatio)
	require.Equal(t, bondedAcc.GetCoins().AmountOf("cet"), bondPool.BondedTokens)
	require.Equal(t, notBondedAcc.GetCoins().AmountOf("cet"), bondPool.NotBondedTokens)
}

func TestCalcBondedRatio(t *testing.T) {
	bondPool := BondPool{
		BondedTokens:      sdk.NewInt(10e8),
		NotBondedTokens:   sdk.NewInt(500e8),
		NonBondableTokens: sdk.NewInt(10000),
		TotalSupply:       sdk.NewInt(510e8),
	}
	expectedBondRatio := bondPool.BondedTokens.ToDec().QuoInt(bondPool.TotalSupply.Sub(bondPool.NonBondableTokens))
	require.Equal(t, expectedBondRatio, calcBondedRatio(&bondPool))
}

func TestCalcBondedRatioNegative(t *testing.T) {
	bondPool := BondPool{
		BondedTokens:      sdk.NewInt(-10e8),
		NotBondedTokens:   sdk.NewInt(500e8),
		NonBondableTokens: sdk.NewInt(10000),
		TotalSupply:       sdk.NewInt(510e8),
	}
	require.Equal(t, sdk.ZeroDec(), calcBondedRatio(&bondPool))
}
