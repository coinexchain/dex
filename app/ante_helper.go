package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/x/authx"
)

var _ authx.AnteHelper = anteHelper{}

type anteHelper struct {
}

func (ah anteHelper) IsMemoRequired(msg sdk.Msg, ctx sdk.Context) bool {
	// TODO
	return false
}

func (ah anteHelper) GasFee(msg sdk.Msg) sdk.Coins {
	// TODO
	return nil
}
