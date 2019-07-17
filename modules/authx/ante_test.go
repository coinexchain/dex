package authx

import (
	"github.com/coinexchain/dex/modules/authx/types"
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
	InitGenesis(testInput.ctx, testInput.axk, DefaultGenesisState())

	ah := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		return ctx, sdk.Result{}, false
	}
	ah2 := wrapAnteHandler(ah, testInput.axk, nil)

	ctx := testInput.ctx.WithBlockHeight(1)
	tx := auth.StdTx{Fee: auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(1))), Gas: 10000000}}
	_, res, abort := ah2(ctx, tx, false)
	require.True(t, abort)
	require.Equal(t, types.CodeSpaceAuthX, res.Codespace)
	require.Equal(t, types.CodeGasPriceTooLow, res.Code)
}

type testAnteHelper struct {
	error sdk.Error
}

func (h testAnteHelper) CheckMsg(ctx sdk.Context, msg sdk.Msg, memo string) sdk.Error {
	return h.error
}

func TestAdditionalError(t *testing.T) {
	ah := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		return ctx, sdk.Result{}, false
	}

	expectedErr := sdk.NewError(sdk.CodespaceRoot, sdk.CodeInternal, "stop here")
	tx := auth.StdTx{Msgs: []sdk.Msg{nil}}
	ah2 := wrapAnteHandler(ah, AccountXKeeper{}, testAnteHelper{error: expectedErr})

	_, res, abort := ah2(sdk.Context{}, tx, true)
	require.True(t, abort)
	require.Equal(t, expectedErr.Result(), res)
}
