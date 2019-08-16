package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
