package govx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type GovBankKeeper struct {
	// The reference to the CoinKeeper to modify balances
	//ck gov.BankKeeper

	ak auth.AccountKeeper

	dk DistributionKeeper
}

func NewKeeper( /*ck gov.BankKeeper, */ ak auth.AccountKeeper, dk DistributionKeeper) GovBankKeeper {

	return GovBankKeeper{
		//ck: ck,
		ak: ak,
		dk: dk,
	}
}

//func (k GovBankKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
//	return k.ck.GetCoins(ctx, addr)
//}

//func (k GovBankKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
//	k.ck.SetSendEnabled(ctx, enabled)
//}

func (k GovBankKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress,
	amt sdk.Coins) /*sdk.Tags, */ sdk.Error {

	//if fromAddr.Equals(gov.DepositedCoinsAccAddr) && toAddr.Equals(gov.BurnedDepositCoinsAccAddr) {
	//	_, err := subtractCoins(ctx, k.ak, fromAddr, amt)
	//	if err != nil {
	//		ctx.Logger().Error("subtractCoins error %v", err)
	//		panic(err)
	//	}
	//
	//	// update FeePool
	//	feePool := k.dk.GetFeePool(ctx)
	//	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(amt))
	//	k.dk.SetFeePool(ctx, feePool)
	//	ctx.Logger().Info("burnt token %v send to community pool", amt)
	//	return nil, nil
	//}
	//
	//return k.ck.SendCoins(ctx, fromAddr, toAddr, amt)
	return nil
}

func subtractCoins(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()

	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := setCoins(ctx, ak, addr, newCoins)

	return newCoins, err
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}
