package keepers

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

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

func (keeper *CommentCountKeeper) IncrCommentCount(ctx sdk.Context, denom string) uint64 {
	store := ctx.KVStore(keeper.commentKey)
	ccKey := append(CommentCountKey, []byte(denom)...)
	a := store.Get(ccKey)
	count := uint64(0)
	if len(a) != 0 {
		count = binary.LittleEndian.Uint64(a)
	}
	lastCount := count
	count++
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, count)
	store.Set(ccKey, b)
	return lastCount
}

func (keeper *CommentCountKeeper) GetAllCommentCount(ctx sdk.Context) map[string]uint64 {
	res := make(map[string]uint64)
	store := ctx.KVStore(keeper.commentKey)
	iter := store.Iterator(CommentCountKey, CommentCountKeyEnd)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := iter.Key()[1:]
		a := iter.Value()
		count := binary.LittleEndian.Uint64(a[:])
		res[string(denom)] = count
	}
	return res
}

func (keeper *CommentCountKeeper) GetCommentCount(ctx sdk.Context, denom string) uint64 {
	store := ctx.KVStore(keeper.commentKey)
	ccKey := append(CommentCountKey, []byte(denom)...)
	a := store.Get(ccKey)
	count := uint64(0)
	if len(a) != 0 {
		count = binary.LittleEndian.Uint64(a[:])
	}
	return count
}

func (keeper *CommentCountKeeper) SetCommentCount(ctx sdk.Context, denom string, count uint64) {
	store := ctx.KVStore(keeper.commentKey)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], count)
	ccKey := append(CommentCountKey, []byte(denom)...)
	store.Set(ccKey, b[:])
}

type Keeper struct {
	cck               *CommentCountKeeper
	bxk               types.ExpectedBankxKeeper
	axk               types.ExpectedAssetStatusKeeper
	ak                types.ExpectedAccountKeeper
	dk                types.ExpectedDistributionxKeeper
	eventTypeMsgQueue string
}

func NewKeeper(key sdk.StoreKey,
	bxk types.ExpectedBankxKeeper,
	axk types.ExpectedAssetStatusKeeper,
	ak types.ExpectedAccountKeeper,
	dk types.ExpectedDistributionxKeeper,
	et string) *Keeper {
	return &Keeper{
		cck:               NewCommentCountKeeper(key),
		bxk:               bxk,
		axk:               axk,
		ak:                ak,
		dk:                dk,
		eventTypeMsgQueue: et,
	}
}

func (k *Keeper) GetEventTypeMsgQueue() string {
	return k.eventTypeMsgQueue
}

func (k *Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.bxk.SendCoins(ctx, from, to, amt)
}

func (k *Keeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return k.axk.IsTokenExists(ctx, denom)
}

func (k *Keeper) DonateToCommunityPool(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.dk.DonateToCommunityPool(ctx, fromAddr, amt)
}

func (k *Keeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account {
	return k.ak.GetAccount(ctx, addr)
}

func (k *Keeper) IncrCommentCount(ctx sdk.Context, denom string) uint64 {
	return k.cck.IncrCommentCount(ctx, denom)
}

func (k *Keeper) GetAllCommentCount(ctx sdk.Context) map[string]uint64 {
	return k.cck.GetAllCommentCount(ctx)
}

func (k *Keeper) GetCommentCount(ctx sdk.Context, denom string) uint64 {
	return k.cck.GetCommentCount(ctx, denom)
}

func (k *Keeper) SetCommentCount(ctx sdk.Context, denom string, count uint64) {
	k.cck.SetCommentCount(ctx, denom, count)
}
