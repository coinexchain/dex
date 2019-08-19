package simulation

import (
	"encoding/binary"

	"math/rand"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/comment/internal/keepers"
	"github.com/coinexchain/dex/modules/comment/internal/types"
	simulationx "github.com/coinexchain/dex/simulation"
)

func randomUTF8OrBytes(r *rand.Rand, length int, isUTF8 bool) []byte {

	var res []byte
	for i := 0; i < length; {
		randomInt := r.Uint32()
		bz := make([]byte, 4)
		binary.LittleEndian.PutUint32(bz, randomInt)
		if isUTF8 && !utf8.Valid(bz) {
			continue
		}
		res = append(res, bz...)
		i = i + len(bz)
	}
	return res

}
func randomToken(r *rand.Rand, ctx sdk.Context, ask asset.Keeper) asset.Token {
	tokenList := ask.GetAllTokens(ctx)
	if len(tokenList) == 0 {
		return nil
	}
	token := tokenList[simulationx.GetRandomElemIndex(r, len(tokenList))]
	return token
}

func randomContent(r *rand.Rand) (contentType int8, content []byte) {

	contentType = int8(r.Intn(int(types.RawBytes + 1)))
	contentLength := r.Intn(types.MaxContentSize)
	switch contentType {
	case types.RawBytes:
		content = randomUTF8OrBytes(r, contentLength, false)
	default:
		content = randomUTF8OrBytes(r, contentLength, true)
	}
	return

}
func randomTokenCommentRef(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, ask asset.Keeper) (token asset.Token, ids []uint64) {

	totalComment := k.Cck.GetAllCommentCount(ctx)
	if len(totalComment) == 0 {
		return nil, []uint64{}
	}
	tokenList := ask.GetAllTokens(ctx)

	for {
		token = tokenList[simulationx.GetRandomElemIndex(r, len(tokenList))]
		total := totalComment[token.GetSymbol()]
		if total == 0 {
			continue
		}
		randLen := r.Intn(int(total))
		for i := 0; i < randLen; i = i + 1 {
			ids = append(ids, uint64(r.Intn(int(total))))
		}
		return
	}
}
