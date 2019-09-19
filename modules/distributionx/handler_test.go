package distributionx_test

import (
	"testing"

	"github.com/coinexchain/dex/modules/distributionx"
	"github.com/coinexchain/dex/testapp"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/coinexchain/dex/modules/bankx"
	types2 "github.com/coinexchain/dex/modules/distributionx/types"
	dex "github.com/coinexchain/dex/types"
)

var validCoins = dex.NewCetCoins(10e8)

type testInput struct {
	k   distributionx.Keeper
	ctx sdk.Context
	ak  auth.AccountKeeper
	dk  distribution.Keeper
}

func setupTestInput() testInput {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.BankKeeper.SetSendEnabled(ctx, true)
	testApp.DistrKeeper.SetFeePool(ctx, distribution.InitialFeePool())

	return testInput{
		ctx: ctx,
		k:   testApp.DistrxKeeper,
		ak:  testApp.AccountKeeper,
		dk:  testApp.DistrKeeper,
	}
}
func TestDonateToCommunityPool(t *testing.T) {
	input := setupTestInput()

	addr := sdk.AccAddress([]byte("addr"))
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := dex.NewCetCoins(1e8)
	msg := types2.NewMsgDonateToCommunityPool(addr, Coins)
	res := distributionx.HandleMsgDonateToCommunityPool(input.ctx, input.k, msg)

	require.True(t, res.IsOK())
	feePool = input.dk.GetFeePool(input.ctx)
	require.True(t, feePool.CommunityPool.AmountOf("cet").Equal(sdk.NewDec(1e8)))

	fromAcc = input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(9e8)))
}

func TestDonateToCommunityPoolFailed(t *testing.T) {
	input := setupTestInput()

	addr := sdk.AccAddress([]byte("addr"))
	acc := auth.NewBaseAccountWithAddress(addr)
	acc.SetCoins(validCoins)
	input.ak.SetAccount(input.ctx, &acc)

	fromAcc := input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

	feePool := input.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	Coins := dex.NewCetCoins(11e8)
	msg := types2.NewMsgDonateToCommunityPool(addr, Coins)
	res := distributionx.HandleMsgDonateToCommunityPool(input.ctx, input.k, msg)

	require.False(t, res.IsOK())
	feePool = input.dk.GetFeePool(input.ctx)
	require.Nil(t, feePool.CommunityPool)

	fromAcc = input.ak.GetAccount(input.ctx, addr)
	require.True(t, fromAcc.GetCoins().AmountOf("cet").Equal(sdk.NewInt(10e8)))

}
func TestNewHandler(t *testing.T) {
	input := setupTestInput()
	handler := distributionx.NewHandler(input.k)

	msg := bankx.MsgSetMemoRequired{}
	res := handler(input.ctx, msg)

	require.Equal(t, sdk.CodeUnknownRequest, res.Code)
}
