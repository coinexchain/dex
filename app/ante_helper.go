package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/stakingx"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
	accountXKeeper authx.AccountXKeeper
}

func (ah anteHelper) CheckMsg(ctx sdk.Context, msg sdk.Msg, memo string) sdk.Error {
	switch msg := msg.(type) {
	case bank.MsgSend: // should not be here!
		return ah.checkMemo(ctx, msg.ToAddress, memo)
	case bankx.MsgSend:
		return ah.checkMemo(ctx, msg.ToAddress, memo)
	case staking.MsgCreateValidator:
		return ah.checkMinSelfDelegation(ctx, msg.MinSelfDelegation)
	}
	return nil
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
	expected := ah.accountXKeeper.GetParams(ctx).MinSelfDelegation
	if actual.LT(expected) {
		return stakingx.ErrMinSelfDelegationBelowRequired(
			fmt.Sprintf("expected:%v actual:%v", expected, actual))
	}
	return nil
}

func (ah anteHelper) GasFee(msg sdk.Msg) sdk.Coins {
	// TODO
	return nil
}
