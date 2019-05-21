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

func (ah anteHelper) CheckMemo(msg sdk.Msg, memo string, ctx sdk.Context) sdk.Error {
	memoLen := len(memo)
	switch msg := msg.(type) {
	case bank.MsgSend:
		if ax, ok := ah.accountXKeeper.GetAccountX(ctx, msg.ToAddress); ok && ax.TransferMemoRequired {
			if memoLen == 0 {
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
