package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/stakingx"
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

func checkAddr(msg sdk.Msg) sdk.Error {
	signers := msg.GetSigners()
	for _, signer := range signers {
		if signer.Equals(incentive.PoolAddr) {
			return sdk.ErrUnauthorized("tx not allowed to be sent from the sender addr")
		}
	}
	return nil
}
