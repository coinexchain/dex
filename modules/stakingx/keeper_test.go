package stakingx

import (
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

	_, _, addr := testutil.KeyPubAddr()
	genesisState := GenesisState{
		Params: Params{
			MinSelfDelegation:    sdk.NewInt(DefaultMinSelfDelegation),
			NonBondableAddresses: []sdk.AccAddress{addr},
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
		MinSelfDelegation:    sdk.ZeroInt(),
		NonBondableAddresses: []sdk.AccAddress{addr},
	}
	sxk.SetParams(ctx, testParam)
	acc := auth.BaseAccount{
		Address: addr,
		Coins:   dex.NewCetCoins(1e8),
	}
	sxk.ak.SetAccount(ctx, &acc)

	pool := staking.Pool{
		BondedTokens:    sdk.NewInt(10e8),
		NotBondedTokens: sdk.NewInt(500e8),
	}
	//sxk.sk.SetPool(ctx, pool)

	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(1000)),
	}
	sxk.dk.SetFeePool(ctx, feePool)

	//test
	expectedNonBondableTokens := feePool.CommunityPool.AmountOf("cet").Add(acc.Coins.AmountOf("cet").ToDec())
	expectedTotalSupply := pool.BondedTokens.Add(pool.NotBondedTokens)
	expectedBondRatio := pool.BondedTokens.ToDec().QuoInt(expectedTotalSupply.Sub(expectedNonBondableTokens.RoundInt()))

	bondPool := sxk.CalcBondPoolStatus(ctx)
	require.Equal(t, pool.NotBondedTokens, bondPool.NotBondedTokens)
	require.Equal(t, pool.BondedTokens, bondPool.BondedTokens)
	require.Equal(t, expectedNonBondableTokens, bondPool.NonBondableTokens.ToDec())
	require.Equal(t, expectedTotalSupply, bondPool.TotalSupply)
	require.Equal(t, expectedBondRatio, bondPool.BondRatio)

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
