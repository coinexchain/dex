package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	dex "github.com/coinexchain/dex/types"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func RandomSymbol(r *rand.Rand, prefix string, randomPartLen int) string {
	bytes := make([]byte, 0, len(prefix)+randomPartLen)
	bytes = append(bytes, []byte(prefix)...)
	for i := 0; i < randomPartLen; i++ {
		bytes = append(bytes, alphabet[r.Intn(36)])
	}
	return string(bytes)
}

func SimulateHandleMsg(msg sdk.Msg, handler sdk.Handler, ctx sdk.Context) (ok bool) {
	ctx, write := ctx.CacheContext()
	ok = handler(ctx, msg).IsOK()
	if ok {
		write()
	}
	return ok
}
func RandomCET(r *rand.Rand, ctx sdk.Context, ak auth.AccountKeeper, fromAcc simulation.Account) (donation int64) {

	fromCoins := ak.GetAccount(ctx, fromAcc.Address).GetCoins().AmountOf(dex.DefaultBondDenom)
	if !fromCoins.IsZero() {
		donation, err := simulation.RandPositiveInt(r, fromCoins)
		if err == nil {
			return donation.Int64()
		}
	}
	return
}
func RandomAccCoins(r *rand.Rand, account auth.Account) (string, int64) {
	coins := account.GetCoins()
	if len(coins) == 0 {
		return "", 0
	}
	randomCoins := coins[r.Intn(len(coins))]
	amt, err := simulation.RandPositiveInt(r, randomCoins.Amount)
	if err == nil {
		return randomCoins.Denom, amt.Int64()
	}
	return "", 0
}
func GetRandomElemIndex(r *rand.Rand, length int) int {
	return r.Intn(length)
}

func RandomBool(r *rand.Rand) bool {
	return r.Uint32()%2 == 0
}
