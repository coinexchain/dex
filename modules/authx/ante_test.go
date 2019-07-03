package authx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func TestOriginalAnteHandlerError(t *testing.T) {
	expectedRes := sdk.NewError(sdk.CodespaceRoot, sdk.CodeInternal, "stop here").Result()
	ah := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		return ctx, expectedRes, true
	}
	ah2 := wrapAnteHandler(ah, AccountXKeeper{}, nil)
	_, res, _ := ah2(sdk.Context{}, nil, false)
	require.Equal(t, expectedRes, res)
}

func TestGasPriceTooLowError(t *testing.T) {
	testInput := setupTestInput()
	testInput.axk.SetParams(testInput.ctx, DefaultParams())

	ah := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		return ctx, sdk.Result{}, false
	}
	ah2 := wrapAnteHandler(ah, testInput.axk, nil)

	ctx := testInput.ctx.WithBlockHeight(1)
	tx := auth.StdTx{Fee: auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(1))), Gas: 10000000}}
	_, res, _ := ah2(ctx, tx, false)
	require.Equal(t, CodeSpaceAuthX, res.Codespace)
	require.Equal(t, CodeGasPriceTooLow, res.Code)
}

func TestAdditionalError(t *testing.T) {
	// TODO
}
