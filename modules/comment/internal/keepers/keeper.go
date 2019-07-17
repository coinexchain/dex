package keepers

import (
	"encoding/binary"
	"github.com/coinexchain/dex/modules/comment/internal/types"

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
	Cck         *CommentCountKeeper
	Bxk         types.ExpectedBankxKeeper
	Axk         types.ExpectedAssetStatusKeeper
	Dk          types.ExpectedDistributionKeeper
	MsgSendFunc func(key string, v interface{}) error
}

func NewKeeper(cck *CommentCountKeeper,
	bxk types.ExpectedBankxKeeper,
	axk types.ExpectedAssetStatusKeeper,
	dk types.ExpectedDistributionKeeper,
	msgSendFunc func(key string, v interface{}) error) *Keeper {
	return &Keeper{
		Cck:         cck,
		Bxk:         bxk,
		Axk:         axk,
		Dk:          dk,
		MsgSendFunc: msgSendFunc,
	}
}
