package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/x/authx"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
	accountXKeeper authx.AccountXKeeper
}

func (ah anteHelper) IsMemoRequired(msg sdk.Msg, ctx sdk.Context) bool {
	switch msg := msg.(type) {
	case bank.MsgSend:
		if ax, ok := ah.accountXKeeper.GetAccountX(ctx, msg.ToAddress); ok {
			return ax.TransferMemoRequired
		}
	}
	return false
}

func (ah anteHelper) GasFee(msg sdk.Msg) sdk.Coins {
	// TODO
	return nil
}
