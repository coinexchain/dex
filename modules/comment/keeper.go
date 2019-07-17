package comment

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	CommentCountKey = []byte{0x10}
)

type CommentCountKeeper struct {
	commentKey sdk.StoreKey
}

func NewCommentCountKeeper(key sdk.StoreKey) *CommentCountKeeper {
	return &CommentCountKeeper{
		commentKey: key,
	}
}

func (keeper *CommentCountKeeper) IncrCommentCount(ctx sdk.Context) {
	store := ctx.KVStore(keeper.commentKey)
	a := store.Get(CommentCountKey)
	count := binary.LittleEndian.Uint64(a[:])
	count++
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], count)
	store.Set(CommentCountKey, b[:])
}

func (keeper *CommentCountKeeper) GetCommentCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(keeper.commentKey)
	a := store.Get(CommentCountKey)
	count := binary.LittleEndian.Uint64(a[:])
	return count
}

func (keeper *CommentCountKeeper) SetCommentCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(keeper.commentKey)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], count)
	store.Set(CommentCountKey, b[:])
}

type Keeper struct {
	cck         *CommentCountKeeper
	bxk         ExpectedBankxKeeper
	axk         ExpectedAssetStatusKeeper
	dk          ExpectedDistributionKeeper
	msgSendFunc func(key string, v interface{}) error
}

func NewKeeper(cck *CommentCountKeeper,
	bxk ExpectedBankxKeeper,
	axk ExpectedAssetStatusKeeper,
	dk ExpectedDistributionKeeper,
	msgSendFunc func(key string, v interface{}) error) *Keeper {
	return &Keeper{
		cck:         cck,
		bxk:         bxk,
		axk:         axk,
		dk:          dk,
		msgSendFunc: msgSendFunc,
	}
}
