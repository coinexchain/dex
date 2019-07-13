package crisisx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	dType "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
)

const (
	ModuleName = "crisisx"
)

func RegisterInvariants(c *crisis.Keeper, tokenKeeper asset.Keeper, bk ExpectBankxKeeper, supplyKeeper ExpectSupplyKeeper, feeCollectorName string, disk distribution.Keeper, stk staking.Keeper) {
	c.RegisterRoute(ModuleName, "cet-invariant", SupplyCETInvariant(tokenKeeper, bk, supplyKeeper, feeCollectorName, disk, stk))
}

func SupplyCETInvariant(tokenKeeper asset.Keeper, bk ExpectBankxKeeper, supplyKeeper ExpectSupplyKeeper, feeCollectorName string, disk distribution.Keeper, stk staking.Keeper) sdk.Invariant {

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
