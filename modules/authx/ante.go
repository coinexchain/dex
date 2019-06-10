package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type AnteHelper interface {
	CheckMsg(ctx sdk.Context, msg sdk.Msg, memo string) sdk.Error
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(ak auth.AccountKeeper, fck auth.FeeCollectionKeeper,
	anteHelper AnteHelper) sdk.AnteHandler {

	ah := auth.NewAnteHandler(ak, fck)
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// run auth.AnteHandler first
		newCtx, res, abort = ah(ctx, tx, simulate)

		// then, do additional check
		stdTx, _ := tx.(auth.StdTx)
		res2 := doAdditionalCheck(ctx, stdTx, anteHelper)
		if !res2.IsOK() {
			res = res2
			abort = true
		}

		return
	}
}

func doAdditionalCheck(ctx sdk.Context, tx auth.StdTx, anteHelper AnteHelper) sdk.Result {
	memo := tx.Memo
	for _, msg := range tx.Msgs {
		if err := anteHelper.CheckMsg(ctx, msg, memo); err != nil {
			return err.Result()
		}
	}
	return sdk.Result{}
}
