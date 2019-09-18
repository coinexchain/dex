package stakingx

import (
	dex "github.com/coinexchain/dex/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func CalcBondedRatio(p *BondPool) sdk.Dec {
	return calcBondedRatio(p)
}

type MockKeeper struct {
	Keeper
	Sk           *staking.Keeper
	Dk           DistributionKeeper
	Ak           auth.AccountKeeper
	Bk           ExpectBankxKeeper
	SupplyKeeper ExpectSupplyKeeper
}

func InitStates(ctx sdk.Context, sxk Keeper, ak auth.AccountKeeper, splk supply.Keeper) MockKeeper {
	//intialize params & states needed
	params := staking.DefaultParams()
	params.BondDenom = "cet"
	sxk.sk.SetParams(ctx, params)

	//initialize FeePool
	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(0)),
	}
	sxk.dk.SetFeePool(ctx, feePool)

	//initialize staking Pool
	bondedAcc := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	notBondedAcc := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	ak.SetAccount(ctx, bondedAcc)
	ak.SetAccount(ctx, notBondedAcc)

	//initialize total supply
	splk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{sdk.NewInt64Coin("cet", 10e8)}})
	return MockKeeper{
		sxk,
		sxk.sk,
		sxk.dk,
		sxk.ak,
		sxk.bk,
		sxk.supplyKeeper,
	}
}
