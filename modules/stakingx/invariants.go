package stakingx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/staking"

	dType "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
)

type AssetViewKeeper interface {
	GetToken(ctx sdk.Context, symbol string) asset.Token
}

func RegisterInvariants(c crisis.Keeper, k Keeper, assetKeeper asset.Keeper, bk ExpectBankxKeeper, supplyKeeper ExpectSupplyKeeper, feeCollectorName string, disk distribution.Keeper, stk staking.Keeper) {

	c.RegisterRoute(ModuleName, "total-supply", TotalSupplyInvariants(k, assetKeeper))

	c.RegisterRoute(ModuleName, "cet-invariant", SupplyCETInvariant(assetKeeper, bk, supplyKeeper, feeCollectorName, disk, stk))

	// SupplyInvariants no longer suitable here, new SupplyInvariants will be created
	// c.RegisterRoute(types.ModuleName, "supply",
	//	SupplyInvariants(k, f, d, am))

	//c.RegisterRoute(types.ModuleName, "nonnegative-power",
	//	staking.NonNegativePowerInvariant(sk))
	//
	//c.RegisterRoute(types.ModuleName, "positive-delegation",
	//	staking.PositiveDelegationInvariant(sk))
	//
	//c.RegisterRoute(types.ModuleName, "delegator-shares",
	//	staking.DelegatorSharesInvariant(sk))
}

func TotalSupplyInvariants(k Keeper, assetKeeper AssetViewKeeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		token := assetKeeper.GetToken(ctx, types.DefaultBondDenom)
		if token == nil {
			return fmt.Errorf("cet not found")
		}

		ts := token.GetTotalSupply()
		bondPool := k.CalcBondPoolStatus(ctx)

		if ts != bondPool.TotalSupply.Int64() {
			return fmt.Errorf("total-supply invariance:\n"+
				"\tinconsistent total-supply: \n"+
				"\tCET asset total supply: %v\n"+
				"\tpool: %v", ts, bondPool)
		}

		return nil
	}
}

func SupplyCETInvariant(tokenKeeper asset.Keeper, bk ExpectBankxKeeper, supplyKeeper ExpectSupplyKeeper,
	feeCollectorName string, disk distribution.Keeper, stk staking.Keeper) sdk.Invariant {

	return func(ctx sdk.Context) error {
		token := tokenKeeper.GetToken(ctx, types.DefaultBondDenom)
		if token == nil {
			return fmt.Errorf("cet not found")
		}

		var totalAmount = sdk.ZeroInt()

		// Get all amounts based on the account system
		basedAccountTotalAmount := bk.TotalAmountOfCoin(ctx, types.CET)
		totalAmount = totalAmount.Add(basedAccountTotalAmount)
		//fmt.Printf("basedAccountTotalAmount : %s, totalAmount : %s \n", basedAccountTotalAmount, totalAmount.String())

		// Get all amounts based on the Non-account system
		feeAmount := GetCollectedFee(ctx, supplyKeeper, feeCollectorName)

		communityAmount := disk.GetFeePool(ctx).CommunityPool.AmountOf(types.CET)
		totalAmount = totalAmount.Add(feeAmount).Add(communityAmount.RoundInt())
		//fmt.Printf("feeAmount : %s, communityAmount : %s, totalAmount : %s\n",
		//	feeAmount.String(), communityAmount.String(), totalAmount.String())

		// Get all amounts based on the Non-account system in the validators
		outStandingProcess := func(val sdk.ValAddress, rewards dType.ValidatorOutstandingRewards) (stop bool) {
			totalAmount = totalAmount.Add(rewards.AmountOf(types.CET).RoundInt())
			//fmt.Printf("validator addr : %s, rewards : %s, totalAmount : %s\n",
			//	val.String(), rewards.AmountOf(types.CET).RoundInt().String(), totalAmount.String())
			return false
		}

		validatorProcess := func(index int64, validator exported.ValidatorI) bool {
			totalAmount = totalAmount.Add(validator.GetTokens())
			//fmt.Printf("validator addr : %s, tokens : %d, totalTokens : %s\n",
			//	validator.GetOperator().String(), validator.GetTokens().Int64(), totalAmount.String())
			return false
		}

		unbondingProcess := func(index int64, ubd staking.UnbondingDelegation) bool {
			for _, ubdentry := range ubd.Entries {
				totalAmount = totalAmount.Add(ubdentry.Balance)
				//fmt.Printf("unbonding tokens : %s,  totalTokens : %s\n, delgator addr : %s, validator addr : %s ",
				//	ubdentry.Balance.String(), totalAmount.String(), ubd.DelegatorAddress.String(), ubd.ValidatorAddress.String())
			}
			return false
		}

		disk.IterateValidatorOutstandingRewards(ctx, outStandingProcess)
		stk.IterateUnbondingDelegations(ctx, unbondingProcess)
		stk.IterateValidators(ctx, validatorProcess)

		// Judge equality
		if totalAmount.Int64() != token.GetTotalSupply() {
			return fmt.Errorf("the cet total amount [ %d ]is inconsistent with the actual amount [ %d ]",
				token.GetTotalSupply(), totalAmount.Int64())
		}

		return nil
	}
}

func GetCollectedFee(ctx sdk.Context, supplyKeeper ExpectSupplyKeeper, feeCollectorName string) sdk.Int {
	feeCollector := supplyKeeper.GetModuleAccount(ctx, feeCollectorName)
	feesCollected := feeCollector.GetCoins()
	return feesCollected.AmountOf(types.CET)
}
