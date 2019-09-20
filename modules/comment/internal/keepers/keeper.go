package keepers

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/comment/internal/types"
)

var (
	CommentCountKey    = []byte{0x10}
	CommentCountKeyEnd = []byte{0x11}
)

type CommentCountKeeper struct {
	commentKey sdk.StoreKey
}

func NewCommentCountKeeper(key sdk.StoreKey) *CommentCountKeeper {
	return &CommentCountKeeper{
		commentKey: key,
	}
}

func (keeper *CommentCountKeeper) IncrCommentCount(ctx sdk.Context, denorm string) uint64 {
	store := ctx.KVStore(keeper.commentKey)
	ccKey := append(CommentCountKey, []byte(denorm)...)
	a := store.Get(ccKey)
	count := uint64(0)
	if len(a) != 0 {
		count = binary.LittleEndian.Uint64(a[:])
	}
	lastCount := count
	count++
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], count)
	store.Set(ccKey, b[:])
	return lastCount
}

func (keeper *CommentCountKeeper) GetAllCommentCount(ctx sdk.Context) map[string]uint64 {
	res := make(map[string]uint64)
	store := ctx.KVStore(keeper.commentKey)
	iter := store.Iterator(CommentCountKey, CommentCountKeyEnd)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denorm := iter.Key()[1:]
		a := iter.Value()
		count := binary.LittleEndian.Uint64(a[:])
		res[string(denorm)] = count
	}
	return res
}

func (keeper *CommentCountKeeper) GetCommentCount(ctx sdk.Context, denorm string) uint64 {
	store := ctx.KVStore(keeper.commentKey)
	ccKey := append(CommentCountKey, []byte(denorm)...)
	a := store.Get(ccKey)
	count := uint64(0)
	if len(a) != 0 {
		count = binary.LittleEndian.Uint64(a[:])
	}
	return count
}

func (keeper *CommentCountKeeper) SetCommentCount(ctx sdk.Context, denorm string, count uint64) {
	store := ctx.KVStore(keeper.commentKey)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], count)
	ccKey := append(CommentCountKey, []byte(denorm)...)
	store.Set(ccKey, b[:])
}

type Keeper struct {
	Cck               *CommentCountKeeper
	Bxk               types.ExpectedBankxKeeper
	Axk               types.ExpectedAssetStatusKeeper
	Ak                types.ExpectedAccountKeeper
	Dk                types.ExpectedDistributionxKeeper
	EventTypeMsgQueue string
}

func NewKeeper(key sdk.StoreKey,
	bxk types.ExpectedBankxKeeper,
	axk types.ExpectedAssetStatusKeeper,
	ak types.ExpectedAccountKeeper,
	dk types.ExpectedDistributionxKeeper,
	et string) *Keeper {
	return &Keeper{
		Cck:               NewCommentCountKeeper(key),
		Bxk:               bxk,
		Axk:               axk,
		Ak:                ak,
		Dk:                dk,
		EventTypeMsgQueue: et,
	}
}
