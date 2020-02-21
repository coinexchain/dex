package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/cet-sdk/types"

	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/bankx"
	"github.com/coinexchain/cet-sdk/modules/distributionx"
	"github.com/coinexchain/cet-sdk/modules/incentive"
	"github.com/coinexchain/cet-sdk/modules/stakingx"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
	accountXKeeper authx.AccountXKeeper
	stakingXKeeper stakingx.Keeper
}

func newAnteHelper(accountXKeeper authx.AccountXKeeper, stakingXKeeper stakingx.Keeper) anteHelper {
	return anteHelper{
		accountXKeeper: accountXKeeper,
		stakingXKeeper: stakingXKeeper,
	}
}

func (ah anteHelper) CheckMsg(ctx sdk.Context, msg sdk.Msg, memo string) sdk.Error {
	if err := checkAddr(msg); err != nil {
		return err
	}

	switch msg := msg.(type) {
	case bankx.MsgSend:
		return ah.checkMemo(ctx, msg.ToAddress, memo)

	case bankx.MsgSupervisedSend:
		return ah.checkMemo(ctx, msg.ToAddress, memo)

	case bankx.MsgMultiSend:
		for _, out := range msg.Outputs {
			if err := ah.checkMemo(ctx, out.Address, memo); err != nil {
				return err
			}
		}
		return nil

	case staking.MsgCreateValidator:
		return ah.checkMsgCreateValidator(ctx, msg)

	case staking.MsgEditValidator:
		return ah.checkMsgEditValidator(ctx, msg.CommissionRate)

	case distribution.MsgSetWithdrawAddress:
		if ah.memoRequired(ctx, msg.WithdrawAddress) {
			return distributionx.ErrMemoRequiredWithdrawAddr(msg.WithdrawAddress.String())
		}
		return nil

	case gov.MsgDeposit:
		return ah.checkMsgDeposit(msg)

	}

	return nil
}

func (ah anteHelper) checkMsgEditValidator(ctx sdk.Context, newRate *sdk.Dec) sdk.Error {
	if newRate == nil {
		return nil
	}

	return ah.checkMinMandatoryCommissionRate(ctx, *newRate)
}

func (ah anteHelper) checkMsgCreateValidator(ctx sdk.Context, msg staking.MsgCreateValidator) sdk.Error {
	if err := ah.checkMinSelfDelegation(ctx, msg.MinSelfDelegation); err != nil {
		return err
	}

	return ah.checkMinMandatoryCommissionRate(ctx, msg.Commission.Rate)
}

func (ah anteHelper) checkMemo(ctx sdk.Context, addr sdk.AccAddress, memo string) sdk.Error {
	if ax, ok := ah.accountXKeeper.GetAccountX(ctx, addr); ok && ax.MemoRequired {
		if len(memo) == 0 {
			return bankx.ErrMemoMissing()
		}
	}
	return nil
}

func (ah anteHelper) memoRequired(ctx sdk.Context, addr sdk.AccAddress) bool {
	if ax, ok := ah.accountXKeeper.GetAccountX(ctx, addr); ok && ax.MemoRequired {
		return true
	}
	return false
}

func (ah anteHelper) checkMinSelfDelegation(ctx sdk.Context, actual sdk.Int) sdk.Error {
	expected := ah.stakingXKeeper.GetParams(ctx).MinSelfDelegation
	if actual.LT(sdk.NewInt(expected)) {
		return stakingx.ErrMinSelfDelegationBelowRequired(expected, actual.Int64())
	}
	return nil
}

func (ah anteHelper) checkMinMandatoryCommissionRate(ctx sdk.Context, actualRate sdk.Dec) sdk.Error {
	minMandatoryRate := ah.stakingXKeeper.GetMinMandatoryCommissionRate(ctx)
	if actualRate.LT(minMandatoryRate) {
		return stakingx.ErrRateBelowMinMandatoryCommissionRate(minMandatoryRate, actualRate)
	}

	return nil
}

func (ah anteHelper) checkMsgDeposit(msg gov.MsgDeposit) sdk.Error {
	if msg.Amount.Len() > 1 || (msg.Amount.Len() == 1 && msg.Amount[0].Denom != types.CET) {
		return sdk.ErrInvalidCoins("tx not allowed to deposit other coins than cet")
	}
	return nil
}

func checkAddr(msg sdk.Msg) sdk.Error {
	signers := msg.GetSigners()
	for _, signer := range signers {
		if signer.Equals(incentive.PoolAddr) {
			return sdk.ErrUnauthorized("tx not allowed to be sent from the sender addr")
		}
	}
	return nil
}
