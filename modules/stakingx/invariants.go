package stakingx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	dType "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/coinexchain/dex/types"
)

func TotalSupplyInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		token := k.assetViewKeeper.GetToken(ctx, types.DefaultBondDenom)
		if token == nil {
			return fmt.Errorf("cet not found")
		}

		ts := token.GetTotalSupply()
		bondPool := k.CalcBondPoolStatus(ctx)

		//TODO: compare with supplyKeeper.TotalSupply
		if ts != bondPool.TotalSupply.Int64() {
			return fmt.Errorf("total-supply invariance:\n"+
				"\tinconsistent total-supply: \n"+
				"\tCET asset total supply: %v\n"+
				"\tpool: %v", ts, bondPool)
		}

		return nil
	}
}

func SupplyCETInvariant(k Keeper) sdk.Invariant {
	assetKeeper := k.assetViewKeeper

	return func(ctx sdk.Context) error {
		token := assetKeeper.GetToken(ctx, types.DefaultBondDenom)
		if token == nil {
			return fmt.Errorf("cet not found")
		}

		var totalAmount = sdk.ZeroInt()

		// Get all amounts based on the account system
		basedAccountTotalAmount := k.bk.TotalAmountOfCoin(ctx, types.CET)
		totalAmount = totalAmount.Add(basedAccountTotalAmount)
		//fmt.Printf("basedAccountTotalAmount : %s, totalAmount : %s \n", basedAccountTotalAmount, totalAmount.String())

		// Get all amounts based on the Non-account system
		feeAmount := GetCollectedFee(ctx, k.supplyKeeper, k.feeCollectorName)

		communityAmount := k.dk.GetFeePool(ctx).CommunityPool.AmountOf(types.CET)
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

		k.dk.IterateValidatorOutstandingRewards(ctx, outStandingProcess)
		k.sk.IterateUnbondingDelegations(ctx, unbondingProcess)
		k.sk.IterateValidators(ctx, validatorProcess)

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
