package authx

import (
	"github.com/coinexchain/dex/types"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestEndBlocker(t *testing.T) {
	input := setupTestInput()

	addr1 := sdk.AccAddress("addr1")
	var accX1 = AccountX{Address: addr1, MemoRequired: false}
	coins := LockedCoins{
		LockedCoin{Coin: types.NewCetCoin(1), UnlockTime: input.ctx.BlockHeader().Time.Unix() - 1},
	}
	accX1.LockedCoins = coins
	input.axk.SetAccountX(input.ctx, accX1)
	acc1 := input.ak.NewAccountWithAddress(input.ctx, addr1)
	coin := types.NewCetCoin(20)
	_ = acc1.SetCoins(sdk.Coins{coin})
	input.ak.SetAccount(input.ctx, acc1)

	addr2 := sdk.AccAddress("addr2")
	var accX2 = AccountX{Address: addr2, MemoRequired: false}
	coins = LockedCoins{
		LockedCoin{Coin: sdk.Coin{Denom: "cet", Amount: sdk.NewInt(1)}, UnlockTime: input.ctx.BlockHeader().Time.Unix() + 1},
	}
	accX2.LockedCoins = coins
	input.axk.SetAccountX(input.ctx, accX2)
	acc2 := input.ak.NewAccountWithAddress(input.ctx, addr2)
	_ = acc2.SetCoins(sdk.Coins{coin})
	input.ak.SetAccount(input.ctx, acc2)

	input.axk.InsertUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()-1, addr1)
	input.axk.InsertUnlockedCoinsQueue(input.ctx, input.ctx.BlockHeader().Time.Unix()+1, addr2)
	EndBlocker(input.ctx, input.axk, input.ak)
	acc1 = input.ak.GetAccount(input.ctx, addr1)
	require.Equal(t, int64(21), acc1.GetCoins().AmountOf("cet").Int64())
	acc2 = input.ak.GetAccount(input.ctx, addr2)
	require.Equal(t, int64(20), acc2.GetCoins().AmountOf("cet").Int64())
}
