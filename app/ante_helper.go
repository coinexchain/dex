package app

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
	accountXKeeper authx.AccountXKeeper
}

func (ah anteHelper) CheckMemo(msg sdk.Msg, memo string, ctx sdk.Context) sdk.Error {
	switch msg := msg.(type) {
	case bankx.MsgSend:
		if ax, ok := ah.accountXKeeper.GetAccountX(ctx, msg.ToAddress); ok && ax.MemoRequired {
			if len(memo) == 0 {
				return bankx.ErrMemoMissing()
			}
		}
	}
	return nil
}

func (ah anteHelper) GasFee(msg sdk.Msg) sdk.Coins {
	// TODO
	return nil
}
