package authx_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	dex "github.com/coinexchain/dex/types"
)

func TestEndBlocker(t *testing.T) {
	input := setupTestInput()

	addr1 := sdk.AccAddress("addr1")
	var accX1 = authx.AccountX{Address: addr1, MemoRequired: false}
	coins := authx.LockedCoins{
		authx.NewLockedCoin("cet", sdk.NewInt(1), input.ctx.BlockHeader().Time.Unix()-1),
	}
	accX1.LockedCoins = coins
	input.axk.SetAccountX(input.ctx, accX1)
	acc1 := input.ak.NewAccountWithAddress(input.ctx, addr1)
	coin := dex.NewCetCoin(20)
	_ = acc1.SetCoins(sdk.Coins{coin})
	input.ak.SetAccount(input.ctx, acc1)
	err := input.tk.IssueToken(input.ctx, "CET token", "cet", sdk.NewInt(20), addr1,
		true, true, false, true, "", "", asset.TestIdentityString)
	require.NoError(t, err)

	addr2 := sdk.AccAddress("addr2")
	var accX2 = authx.AccountX{Address: addr2, MemoRequired: false}
	coins = authx.LockedCoins{
		authx.NewLockedCoin("cet", sdk.NewInt(1), input.ctx.BlockHeader().Time.Unix()+1),
	}
	accX2.LockedCoins = coins
	input.axk.SetAccountX(input.ctx, accX2)
	acc2 := input.ak.NewAccountWithAddress(input.ctx, addr2)
	_ = acc2.SetCoins(sdk.Coins{dex.NewCetCoin(20)})
	input.ak.SetAccount(input.ctx, acc2)

	//set module account for authx
	moduleAccount := input.sk.GetModuleAccount(input.ctx, authx.ModuleName)
	res := moduleAccount.SetCoins(dex.NewCetCoins(2))
	require.NoError(t, res)
	input.sk.SetModuleAccount(input.ctx, moduleAccount)

	input.axk.InsertUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()-1, addr1)
	input.axk.RemoveFromUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()-1, addr1)
	input.axk.EventTypeMsgQueue = "test"
	authx.EndBlocker(input.ctx, input.axk, input.ak, input.tk)
	acc1 = input.ak.GetAccount(input.ctx, addr1)
	require.Equal(t, int64(20), acc1.GetCoins().AmountOf("cet").Int64())

	input.axk.InsertUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()-1, addr1)
	input.axk.InsertUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()+1, addr2)
	authx.EndBlocker(input.ctx, input.axk, input.ak, input.tk)
	acc1 = input.ak.GetAccount(input.ctx, addr1)
	require.Equal(t, int64(21), acc1.GetCoins().AmountOf("cet").Int64())
	acc2 = input.ak.GetAccount(input.ctx, addr2)
	require.Equal(t, int64(20), acc2.GetCoins().AmountOf("cet").Int64())
}
