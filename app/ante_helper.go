package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
	accountXKeeper authx.AccountXKeeper
}

func (ah anteHelper) CheckMsg(msg sdk.Msg, memo string, ctx sdk.Context) sdk.Error {
	switch msg := msg.(type) {
	case bank.MsgSend: // should not be here!
		return ah.checkMemo(msg.ToAddress, memo, ctx)
	case bankx.MsgSend:
		return ah.checkMemo(msg.ToAddress, memo, ctx)
	}
	return nil
}

func (ah anteHelper) checkMemo(addr sdk.AccAddress, memo string, ctx sdk.Context) sdk.Error {
	if ax, ok := ah.accountXKeeper.GetAccountX(ctx, addr); ok && ax.MemoRequired {
		if len(memo) == 0 {
			return bankx.ErrMemoMissing()
		}
	}
	return nil
}

func (ah anteHelper) GasFee(msg sdk.Msg) sdk.Coins {
	// TODO
	return nil
}
